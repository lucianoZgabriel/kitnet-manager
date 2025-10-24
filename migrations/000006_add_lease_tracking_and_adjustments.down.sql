-- Migration DOWN: Reverter adição de lease tracking e rent adjustments

-- Remover tabela de reajustes
DROP TABLE IF EXISTS lease_rent_adjustments;

-- Remover campos de rastreamento da tabela leases
ALTER TABLE leases DROP COLUMN IF EXISTS generation;
ALTER TABLE leases DROP COLUMN IF EXISTS parent_lease_id;
