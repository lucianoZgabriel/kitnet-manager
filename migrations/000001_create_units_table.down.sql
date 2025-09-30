-- Drop trigger and function
DROP TRIGGER IF EXISTS update_units_updated_at ON units;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_units_is_renovated;
DROP INDEX IF EXISTS idx_units_floor;
DROP INDEX IF EXISTS idx_units_status;

-- Drop table
DROP TABLE IF EXISTS units;

-- Drop enum type
DROP TYPE IF EXISTS unit_status;