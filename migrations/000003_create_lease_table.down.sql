-- Reverter migration: Drop leases table

-- Remover índices primeiro 
DROP INDEX IF EXISTS idx_leases_tenant_status;
DROP INDEX IF EXISTS idx_leases_unit_status;
DROP INDEX IF EXISTS idx_leases_end_date;
DROP INDEX IF EXISTS idx_leases_status;
DROP INDEX IF EXISTS idx_leases_tenant_id;
DROP INDEX IF EXISTS idx_leases_unit_id;

-- Remover tabela
DROP TABLE IF EXISTS leases;