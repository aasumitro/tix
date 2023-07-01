CREATE TABLE IF NOT EXISTS participants (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    event_id BIGINT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    job VARCHAR(255) NOT NULL,
    pop VARCHAR(255) NOT NULL,
    dob VARCHAR(255) NOT NULL,
    approved_at BIGINT,
    declined_at BIGINT,
    declined_reason TEXT,
    created_at BIGINT NOT NULL DEFAULT extract(epoch from now()),
    updated_at BIGINT
);

-- ALTER TABLE participants
--     ADD CONSTRAINT fk_event_participants
--         FOREIGN KEY (event_id)
--             REFERENCES events(id);