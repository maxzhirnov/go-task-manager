-- Create tasks table if it doesn't exist
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Insert some initial data
INSERT INTO tasks (title, description, status, created_at, updated_at)
VALUES 
    ('Complete project setup', 'Set up the initial project structure and database', 'pending', NOW(), NOW()),
    ('Write documentation', 'Document the API endpoints and setup instructions', 'pending', NOW(), NOW()),
    ('Add unit tests', 'Implement comprehensive unit tests for the project', 'pending', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Create the users table if it doesn't exist
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL, -- Hashed password
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Add a sample user for testing (password: "password123" hashed with bcrypt)
INSERT INTO users (username, password)
VALUES ('admin', '$2a$10$7Q9Q7pG6Mv7C1vH8vP6x8u9X7w7F1HS2j3dJ3dK2a3b4c5d6e7f8g')
ON CONFLICT DO NOTHING;