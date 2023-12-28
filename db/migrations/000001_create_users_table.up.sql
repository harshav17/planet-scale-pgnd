CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(255) PRIMARY KEY,  -- Unique identifier from Auth0
    email VARCHAR(50),
    name VARCHAR(50),
    -- Other application-specific fields
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
