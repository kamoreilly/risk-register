-- Create enums for incidents
CREATE TYPE incident_status AS ENUM ('new', 'acknowledged', 'in_progress', 'on_hold', 'resolved', 'closed');
CREATE TYPE incident_priority AS ENUM ('p1', 'p2', 'p3', 'p4');

-- Create incidents table
CREATE TABLE incidents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category_id UUID REFERENCES incident_categories(id) ON DELETE SET NULL,
    priority incident_priority NOT NULL DEFAULT 'p3',
    status incident_status NOT NULL DEFAULT 'new',
    assignee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    reporter_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    service_affected VARCHAR(255),
    root_cause TEXT,
    resolution_notes TEXT,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    detected_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    updated_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT
);

-- Create incident_risks junction table
CREATE TABLE incident_risks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    risk_id UUID NOT NULL REFERENCES risks(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    UNIQUE(incident_id, risk_id)
);

-- Create indexes for incidents table
CREATE INDEX idx_incidents_status ON incidents(status);
CREATE INDEX idx_incidents_priority ON incidents(priority);
CREATE INDEX idx_incidents_category ON incidents(category_id);
CREATE INDEX idx_incidents_assignee ON incidents(assignee_id);
CREATE INDEX idx_incidents_reporter ON incidents(reporter_id);
CREATE INDEX idx_incidents_occurred_at ON incidents(occurred_at);
CREATE INDEX idx_incidents_resolved_at ON incidents(resolved_at);

-- Create indexes for incident_risks table
CREATE INDEX idx_incident_risks_incident ON incident_risks(incident_id);
CREATE INDEX idx_incident_risks_risk ON incident_risks(risk_id);
