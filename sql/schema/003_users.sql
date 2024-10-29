-- +goose up
ALTER TABLE users ADD COLUMN hashed_password TEXT;
UPDATE users SET hashed_password = 'unset' WHERE hashed_password IS NULL;
ALTER TABLE users ALTER COLUMN hashed_password SET NOT NULL;

-- +goose down
ALTER TABLE users DROP COLUMN hashed_password;