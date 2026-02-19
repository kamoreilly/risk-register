CREATE TYPE mitigation_status AS ENUM ('planned', 'in_progress', 'completed', 'cancelled');

CREATE TABLE mitigations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    risk_id UUID NOT NULL REFERENCES risks(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    owner VARCHAR(255),
    status mitigation_status NOT NULL DEFAULT 'planned',
    due_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_mitigations_risk ON mitigations(risk_id);
