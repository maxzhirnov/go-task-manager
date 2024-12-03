ALTER TABLE users 
    ADD COLUMN IF NOT EXISTS email VARCHAR(255) UNIQUE;

-- Update existing users with temporary emails based on their usernames
UPDATE users 
SET email = CONCAT(username, '@temporary.com')
WHERE email IS NULL;

-- Make email NOT NULL after setting temporary values
ALTER TABLE users 
    ALTER COLUMN email SET NOT NULL;