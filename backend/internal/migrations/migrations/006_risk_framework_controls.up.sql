CREATE TABLE risk_framework_controls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    risk_id UUID NOT NULL REFERENCES risks(id) ON DELETE CASCADE,
    framework_id UUID NOT NULL REFERENCES frameworks(id) ON DELETE CASCADE,
    control_ref VARCHAR(100) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(risk_id, framework_id, control_ref)
);

CREATE INDEX idx_risk_framework_controls_risk ON risk_framework_controls(risk_id);
CREATE INDEX idx_risk_framework_controls_framework ON risk_framework_controls(framework_id);
