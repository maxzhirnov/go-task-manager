DROP INDEX IF EXISTS idx_users_reset_token;
ALTER TABLE users
    DROP COLUMN IF EXISTS reset_password_token,
    DROP COLUMN IF EXISTS reset_token_expires;