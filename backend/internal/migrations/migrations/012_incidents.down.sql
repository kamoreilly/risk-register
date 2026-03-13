-- Drop indexes for incident_risks table
DROP INDEX IF EXISTS idx_incident_risks_risk;
DROP INDEX IF EXISTS idx_incident_risks_incident;

-- Drop indexes for incidents table
DROP INDEX IF EXISTS idx_incidents_resolved_at;
DROP INDEX IF EXISTS idx_incidents_occurred_at;
DROP INDEX IF EXISTS idx_incidents_reporter;
DROP INDEX IF EXISTS idx_incidents_assignee;
DROP INDEX IF EXISTS idx_incidents_category;
DROP INDEX IF EXISTS idx_incidents_priority;
DROP INDEX IF EXISTS idx_incidents_status;

-- Drop tables
DROP TABLE IF EXISTS incident_risks;
DROP TABLE IF EXISTS incidents;

-- Drop enums
DROP TYPE IF EXISTS incident_priority;
DROP TYPE IF EXISTS incident_status;
