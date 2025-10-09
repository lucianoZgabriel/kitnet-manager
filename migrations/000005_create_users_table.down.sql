-- Drop trigger
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_username;

-- Drop table
DROP TABLE IF EXISTS users;

-- Drop enum type
DROP TYPE IF EXISTS user_role;