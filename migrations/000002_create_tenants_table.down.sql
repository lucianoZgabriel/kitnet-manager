-- Drop trigger
DROP TRIGGER IF EXISTS update_tenants_updated_at ON tenants;

-- Drop indexes
DROP INDEX IF EXISTS idx_tenants_created_at;
DROP INDEX IF EXISTS idx_tenants_full_name;
DROP INDEX IF EXISTS idx_tenants_cpf;

-- Drop table
DROP TABLE IF EXISTS tenants;