-- 000004_add_email_field.down.sql

-- First remove the NOT NULL constraint
ALTER TABLE users 
    ALTER COLUMN email DROP NOT NULL;

-- Then drop the email column
ALTER TABLE users 
    DROP COLUMN IF EXISTS email;