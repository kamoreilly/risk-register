DROP INDEX IF EXISTS idx_risk_framework_controls_risk_control;
DROP INDEX IF EXISTS idx_risk_framework_controls_framework_control;

ALTER TABLE risk_framework_controls
    ADD COLUMN framework_id UUID REFERENCES frameworks(id) ON DELETE CASCADE,
    ADD COLUMN control_ref VARCHAR(100);

UPDATE risk_framework_controls rfc
SET framework_id = fc.framework_id,
    control_ref = fc.control_ref
FROM framework_controls fc
WHERE fc.id = rfc.framework_control_id;

ALTER TABLE risk_framework_controls
    ALTER COLUMN framework_id SET NOT NULL,
    ALTER COLUMN control_ref SET NOT NULL;

ALTER TABLE risk_framework_controls
    DROP CONSTRAINT IF EXISTS risk_framework_controls_framework_control_id_fkey;

ALTER TABLE risk_framework_controls
    DROP COLUMN framework_control_id;

CREATE INDEX idx_risk_framework_controls_framework ON risk_framework_controls(framework_id);

ALTER TABLE risk_framework_controls
    ADD CONSTRAINT risk_framework_controls_risk_id_framework_id_control_ref_key
        UNIQUE (risk_id, framework_id, control_ref);

DROP INDEX IF EXISTS idx_framework_controls_ref;
DROP INDEX IF EXISTS idx_framework_controls_framework;
DROP TABLE IF EXISTS framework_controls;