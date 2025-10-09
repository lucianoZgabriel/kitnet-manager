-- Create ENUM type for user role
CREATE TYPE user_role AS ENUM (
    'admin',
    'manager',
    'viewer'
);

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'viewer',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT check_username_length CHECK (length(trim(username)) >= 3),
    CONSTRAINT check_password_hash_not_empty CHECK (length(trim(password_hash)) > 0)
);

-- Create indexes for common queries
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);

-- Create trigger to auto-update updated_at
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE users IS 'System users with authentication and authorization';
COMMENT ON COLUMN users.username IS 'Unique username for login';
COMMENT ON COLUMN users.password_hash IS 'Bcrypt hashed password';
COMMENT ON COLUMN users.role IS 'User role: admin (full access), manager (manage operations), viewer (read-only)';
COMMENT ON COLUMN users.is_active IS 'Whether user account is active';
COMMENT ON COLUMN users.last_login_at IS 'Timestamp of last successful login';

-- Insert default admin user
-- Password: admin123 (bcrypt hash)
-- IMPORTANTE: Trocar a senha após primeiro login!
INSERT INTO users (username, password_hash, role) VALUES 
    ('admin', '$2a$10$rQ3Kj7qX5gxV.1Y7hN5qUO5vW0H2KXY8qZ3nJ9wN5xM4pL7kJ6vQi', 'admin');

COMMENT ON TABLE users IS 'Usuário admin padrão criado com senha: admin123';