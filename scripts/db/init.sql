-- Create tasks table if it doesn't exist
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
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