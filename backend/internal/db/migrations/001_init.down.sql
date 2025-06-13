-- +goose Down
-- Revert initial schema for GoChop backend

DROP TABLE IF EXISTS analytics;
DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS verification_tokens;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users;
-- Optionally disable pgcrypto extension (commented):
-- DROP EXTENSION IF EXISTS "pgcrypto"; 