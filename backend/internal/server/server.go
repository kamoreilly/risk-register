package server

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"

	"backend/internal/database"
	"backend/internal/handlers"
)

type FiberServer struct {
	*fiber.App
	db                database.Service
	rawDB             *sql.DB
	users             database.UserRepository
	risks             database.RiskRepository
	categories        database.CategoryRepository
	mitigations       database.MitigationRepository
	frameworks        database.FrameworkRepository
	controls          database.RiskFrameworkControlRepository
	audit             database.AuditLogRepository
	auth              *handlers.AuthHandler
	riskHandler       *handlers.RiskHandler
	categoryHandler   *handlers.CategoryHandler
	mitigationHandler *handlers.MitigationHandler
	frameworkHandler  *handlers.FrameworkHandler
	controlHandler    *handlers.ControlHandler
	dashboardHandler  *handlers.DashboardHandler
	aiHandler         *handlers.AIHandler
	auditHandler      *handlers.AuditHandler
}

func New() *FiberServer {
	db := database.New()
	rawDB := getRawDB()
	users := database.NewUserRepository(rawDB)
	risks := database.NewRiskRepository(rawDB)
	categories := database.NewCategoryRepository(rawDB)
	mitigations := database.NewMitigationRepository(rawDB)
	frameworks := database.NewFrameworkRepository(rawDB)
	controls := database.NewRiskFrameworkControlRepository(rawDB)
	audit := database.NewAuditLogRepository(rawDB)
	dashboard := database.NewDashboardRepository(rawDB)

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "risk-register",
			AppName:      "Risk Register API",
		}),
		db:                db,
		rawDB:             rawDB,
		users:             users,
		risks:             risks,
		categories:        categories,
		mitigations:       mitigations,
		frameworks:        frameworks,
		controls:          controls,
		audit:             audit,
		auth:              handlers.NewAuthHandler(users),
		riskHandler:       handlers.NewRiskHandler(risks, categories, audit),
		categoryHandler:   handlers.NewCategoryHandler(categories),
		mitigationHandler: handlers.NewMitigationHandler(mitigations),
		frameworkHandler:  handlers.NewFrameworkHandler(frameworks, controls),
		controlHandler:    handlers.NewControlHandler(controls),
		dashboardHandler:  handlers.NewDashboardHandler(dashboard),
		aiHandler:         handlers.NewAIHandler(),
		auditHandler:      handlers.NewAuditHandler(audit),
	}

	return server
}

func getRawDB() *sql.DB {
	connStr := buildConnStr()
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	return db
}

func buildConnStr() string {
	host := getEnv("RISK_REGISTER_DB_HOST", "localhost")
	port := getEnv("RISK_REGISTER_DB_PORT", "5432")
	user := getEnv("RISK_REGISTER_DB_USERNAME", "postgres")
	password := getEnv("RISK_REGISTER_DB_PASSWORD", "postgres")
	database := getEnv("RISK_REGISTER_DB_DATABASE", "risk_register")
	schema := getEnv("RISK_REGISTER_DB_SCHEMA", "public")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		user, password, host, port, database, schema)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
