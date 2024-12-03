-- Remove the index first
DROP INDEX IF EXISTS idx_users_username;

-- Re-add the unique constraint
ALTER TABLE users 
    ADD CONSTRAINT users_username_key UNIQUE (username);