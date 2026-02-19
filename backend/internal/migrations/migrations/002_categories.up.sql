CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

INSERT INTO categories (name, description) VALUES
    ('Security', 'Security-related risks including cyber threats and data breaches'),
    ('Operational', 'Operational risks affecting business processes'),
    ('Financial', 'Financial risks including market, credit, and liquidity risks'),
    ('Compliance', 'Regulatory and compliance risks'),
    ('Strategic', 'Strategic risks affecting long-term business objectives'),
    ('Reputational', 'Risks to company reputation and brand');
