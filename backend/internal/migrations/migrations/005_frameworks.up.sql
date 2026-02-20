CREATE TABLE frameworks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_frameworks_name ON frameworks(name);

-- Insert default frameworks
INSERT INTO frameworks (name, description) VALUES
    ('ISO 27001', 'Information Security Management System'),
    ('SOC 2', 'Service Organization Control 2'),
    ('NIST CSF', 'NIST Cybersecurity Framework'),
    ('GDPR', 'General Data Protection Regulation'),
    ('HIPAA', 'Health Insurance Portability and Accountability Act');
