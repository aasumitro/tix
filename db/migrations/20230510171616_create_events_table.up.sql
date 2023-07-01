CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    google_form_id VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    preregister_date BIGINT NOT NULL,
    event_date BIGINT NOT NULL,
    location VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
    updated_at BIGINT
);