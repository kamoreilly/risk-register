DELETE FROM incident_categories WHERE name IN ('Outage', 'Breach', 'Error', 'External', 'Performance');
DROP TABLE IF EXISTS incident_categories;
