CREATE TYPE risk_status AS ENUM ('open', 'mitigating', 'resolved', 'accepted');
CREATE TYPE risk_severity AS ENUM ('low', 'medium', 'high', 'critical');

CREATE TABLE risks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status risk_status NOT NULL DEFAULT 'open',
    severity risk_severity NOT NULL DEFAULT 'medium',
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    review_date DATE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    updated_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT
);

CREATE INDEX idx_risks_owner ON risks(owner_id);
CREATE INDEX idx_risks_status ON risks(status);
CREATE INDEX idx_risks_severity ON risks(severity);
CREATE INDEX idx_risks_category ON risks(category_id);
CREATE INDEX idx_risks_review_date ON risks(review_date);
