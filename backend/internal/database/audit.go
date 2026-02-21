package database

import (
	"context"
	"database/sql"
	"encoding/json"

	"backend/internal/models"

	"github.com/google/uuid"
)

type AuditLogRepository interface {
	Create(ctx context.Context, entityType, entityID string, action models.AuditAction, changes map[string]any, userID string) error
	ListByEntity(ctx context.Context, entityType, entityID string, limit int) ([]*models.AuditLog, error)
}

type auditLogRepo struct {
	db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) AuditLogRepository {
	return &auditLogRepo{db: db}
}

func (r *auditLogRepo) Create(ctx context.Context, entityType, entityID string, action models.AuditAction, changes map[string]any, userID string) error {
	var changesJSON []byte
	var err error
	if changes != nil {
		changesJSON, err = json.Marshal(changes)
		if err != nil {
			return err
		}
	}

	_, err = r.db.ExecContext(ctx,
		`INSERT INTO audit_logs (id, entity_type, entity_id, action, changes, user_id) VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New().String(), entityType, entityID, action, changesJSON, userID,
	)
	return err
}

func (r *auditLogRepo) ListByEntity(ctx context.Context, entityType, entityID string, limit int) ([]*models.AuditLog, error) {
	if limit == 0 {
		limit = 50
	}

	query := `
		SELECT al.id, al.entity_type, al.entity_id, al.action, al.changes, al.user_id, u.name, al.created_at
		FROM audit_logs al
		LEFT JOIN users u ON u.id = al.user_id
		WHERE al.entity_type = $1 AND al.entity_id = $2
		ORDER BY al.created_at DESC
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, entityType, entityID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var changesJSON []byte
		var userName sql.NullString
		if err := rows.Scan(&log.ID, &log.EntityType, &log.EntityID, &log.Action, &changesJSON, &log.UserID, &userName, &log.CreatedAt); err != nil {
			return nil, err
		}
		if userName.Valid {
			log.UserName = userName.String
		}
		if changesJSON != nil {
			json.Unmarshal(changesJSON, &log.Changes)
		}
		logs = append(logs, &log)
	}
	return logs, nil
}
