BEGIN;
--     ALTER TABLE participants DROP CONSTRAINT fk_event_participants;
    DROP TABLE IF EXISTS participants;
COMMIT;
