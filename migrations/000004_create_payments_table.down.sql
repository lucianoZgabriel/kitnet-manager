-- Drop indices primeiro
DROP INDEX IF EXISTS idx_payments_lease_id;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_due_date;
DROP INDEX IF EXISTS idx_payments_payment_type;
DROP INDEX IF EXISTS idx_payments_status_due_date;
DROP INDEX IF EXISTS idx_payments_lease_status;

-- Reverter criação da tabela payments
DROP TABLE IF EXISTS payments CASCADE;