BEGIN;

CREATE TABLE IF NOT EXISTS registrations (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(13) NOT NULL,
    org_name VARCHAR(255),
    designation VARCHAR(255),
    mkt_source VARCHAR(255),
    food_pref VARCHAR(100) NOT NULL,
    t_shirt VARCHAR(10) NOT NULL,
    created_on TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_registrations_created_on ON registrations(created_on);

COMMIT;
