-- Database initialization script
-- This runs when the PostgreSQL container starts for the first time

-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tables are created automatically by GORM, but we can add any custom setup here

-- Example: Create a default admin user (password: admin123)
-- Uncomment the lines below if you want a default admin user
-- INSERT INTO users (email, password, name, created_at, updated_at)
-- VALUES (
--     'admin@example.com',
--     '$2a$12$N9ylo8s5m0e6rC5Vkj8gS.zVVqjC5cZ8iHZL8p0l7k5bF6z1m9X3e',
--     'Admin User',
--     NOW(),
--     NOW()
-- );

-- Example: Create indexes for better query performance
-- CREATE INDEX idx_users_email ON users(email);
-- CREATE INDEX idx_items_user_id ON items(user_id);
-- CREATE INDEX idx_api_keys_key ON api_keys(key);
