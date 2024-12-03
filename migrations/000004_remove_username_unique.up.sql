-- Drop the unique constraint from username
ALTER TABLE users 
    DROP CONSTRAINT IF EXISTS users_username_key;

-- Add an index on username for faster lookups (optional, but recommended)
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);