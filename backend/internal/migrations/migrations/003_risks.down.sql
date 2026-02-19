DROP INDEX IF EXISTS idx_risks_review_date;
DROP INDEX IF EXISTS idx_risks_category;
DROP INDEX IF EXISTS idx_risks_severity;
DROP INDEX IF EXISTS idx_risks_status;
DROP INDEX IF EXISTS idx_risks_owner;
DROP TABLE IF EXISTS risks;
DROP TYPE IF EXISTS risk_severity;
DROP TYPE IF EXISTS risk_status;
