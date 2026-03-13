CREATE TABLE incident_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_incident_categories_name ON incident_categories(name);

INSERT INTO incident_categories (name, description) VALUES
    ('Outage', 'Service or system unavailability incidents'),
    ('Breach', 'Security breach or data compromise incidents'),
    ('Error', 'Application or system error incidents'),
    ('External', 'Incidents caused by external factors or third parties'),
    ('Performance', 'Performance degradation or slowness incidents');
