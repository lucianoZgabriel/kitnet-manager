-- name: GetOccupancyMetrics :one
SELECT
    COUNT(*) as total_units,
    COUNT(*) FILTER (WHERE status = 'occupied') as occupied_units,
    COUNT(*) FILTER (WHERE status = 'available') as available_units,
    COUNT(*) FILTER (WHERE status = 'maintenance') as maintenance_units,
    COUNT(*) FILTER (WHERE status = 'renovation') as renovation_units
FROM units;

-- name: GetMonthlyProjectedRevenue :one
SELECT COALESCE(SUM(monthly_rent_value), 0)::TEXT as total
FROM leases
WHERE status = 'active';

-- name: GetMonthlyRealizedRevenue :one
SELECT COALESCE(SUM(amount), 0)::TEXT as total
FROM payments
WHERE status = 'paid'
  AND payment_type = 'rent'
  AND DATE_TRUNC('month', payment_date) = DATE_TRUNC('month', CURRENT_DATE);

-- name: GetOverdueAmount :one
SELECT COALESCE(SUM(amount), 0)::TEXT as total
FROM payments
WHERE status = 'overdue';

-- name: GetTotalPendingAmount :one
SELECT COALESCE(SUM(amount), 0)::TEXT as total
FROM payments
WHERE status IN ('pending', 'overdue');