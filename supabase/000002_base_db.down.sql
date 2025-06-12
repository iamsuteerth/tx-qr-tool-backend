BEGIN;

DROP INDEX IF EXISTS idx_registrations_created_on;
DROP TABLE IF EXISTS registrations;

COMMIT;