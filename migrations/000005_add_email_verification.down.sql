-- Drop indexes
DROP INDEX IF EXISTS idx_verification_tokens_token;
DROP INDEX IF EXISTS idx_verification_tokens_user_id;
DROP INDEX IF EXISTS idx_verification_tokens_expires_at;

-- Drop verification tokens table
DROP TABLE IF EXISTS verification_tokens;

-- Remove verification status from users
ALTER TABLE users 
    DROP COLUMN IF EXISTS is_verified;