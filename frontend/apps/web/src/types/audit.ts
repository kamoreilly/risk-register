export type AuditAction = 'created' | 'updated' | 'deleted';

export interface AuditLog {
  id: string;
  entity_type: string;
  entity_id: string;
  action: AuditAction;
  changes?: Record<string, { from?: unknown; to?: unknown } | unknown>;
  user_id: string;
  user_name?: string;
  created_at: string;
}
