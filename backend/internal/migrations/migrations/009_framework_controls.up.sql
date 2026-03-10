CREATE TABLE framework_controls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    framework_id UUID NOT NULL REFERENCES frameworks(id) ON DELETE CASCADE,
    control_ref VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(framework_id, control_ref)
);

CREATE INDEX idx_framework_controls_framework ON framework_controls(framework_id);
CREATE INDEX idx_framework_controls_ref ON framework_controls(control_ref);

INSERT INTO framework_controls (framework_id, control_ref, title, description)
SELECT DISTINCT
    rfc.framework_id,
    rfc.control_ref,
    rfc.control_ref,
    NULL
FROM risk_framework_controls rfc;

ALTER TABLE risk_framework_controls
    ADD COLUMN framework_control_id UUID;

UPDATE risk_framework_controls rfc
SET framework_control_id = fc.id
FROM framework_controls fc
WHERE fc.framework_id = rfc.framework_id
  AND fc.control_ref = rfc.control_ref;

ALTER TABLE risk_framework_controls
    ALTER COLUMN framework_control_id SET NOT NULL;

ALTER TABLE risk_framework_controls
    ADD CONSTRAINT risk_framework_controls_framework_control_id_fkey
        FOREIGN KEY (framework_control_id) REFERENCES framework_controls(id) ON DELETE RESTRICT;

ALTER TABLE risk_framework_controls
    DROP CONSTRAINT IF EXISTS risk_framework_controls_risk_id_framework_id_control_ref_key;

DROP INDEX IF EXISTS idx_risk_framework_controls_framework;

ALTER TABLE risk_framework_controls
    DROP COLUMN framework_id,
    DROP COLUMN control_ref;

CREATE INDEX idx_risk_framework_controls_framework_control ON risk_framework_controls(framework_control_id);
CREATE UNIQUE INDEX idx_risk_framework_controls_risk_control ON risk_framework_controls(risk_id, framework_control_id);