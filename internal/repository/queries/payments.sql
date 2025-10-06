-- name: CreatePayment :one
INSERT INTO payments (
    id,
    lease_id,
    payment_type,
    reference_month,
    amount,
    status,
    due_date,
    payment_date,
    payment_method,
    proof_url,
    notes,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING *;

-- name: GetPaymentByID :one
SELECT * FROM payments
WHERE id = $1
LIMIT 1;

-- name: ListPayments :many
SELECT * FROM payments
ORDER BY due_date DESC;

-- name: ListPaymentsByLeaseID :many
SELECT * FROM payments
WHERE lease_id = $1
ORDER BY due_date ASC;

-- name: ListPaymentsByStatus :many
SELECT * FROM payments
WHERE status = $1
ORDER BY due_date ASC;

-- name: GetOverduePayments :many
SELECT * FROM payments
WHERE status IN ('pending', 'overdue')
  AND due_date < CURRENT_DATE
ORDER BY due_date ASC;

-- name: GetUpcomingPayments :many
SELECT * FROM payments
WHERE status = 'pending'
  AND due_date >= CURRENT_DATE
  AND due_date <= CURRENT_DATE + $1::INTEGER
ORDER BY due_date ASC;

-- name: UpdatePayment :one
UPDATE payments
SET
    lease_id = $2,
    payment_type = $3,
    reference_month = $4,
    amount = $5,
    status = $6,
    due_date = $7,
    payment_date = $8,
    payment_method = $9,
    proof_url = $10,
    notes = $11,
    updated_at = $12
WHERE id = $1
RETURNING *;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET
    status = $2,
    updated_at = $3
WHERE id = $1
RETURNING *;

-- name: MarkPaymentAsPaid :one
UPDATE payments
SET
    status = 'paid',
    payment_date = $2,
    payment_method = $3,
    updated_at = $4
WHERE id = $1
RETURNING *;

-- name: MarkPaymentsAsOverdue :exec
UPDATE payments
SET
    status = 'overdue',
    updated_at = NOW()
WHERE status = 'pending'
  AND due_date < CURRENT_DATE;

-- name: CancelPayment :one
UPDATE payments
SET
    status = 'cancelled',
    updated_at = $2
WHERE id = $1
RETURNING *;

-- name: DeletePayment :exec
DELETE FROM payments
WHERE id = $1;

-- name: CountPayments :one
SELECT COUNT(*) FROM payments;

-- name: CountPaymentsByStatus :one
SELECT COUNT(*) FROM payments
WHERE status = $1;

-- name: CountPaymentsByLeaseID :one
SELECT COUNT(*) FROM payments
WHERE lease_id = $1;

-- name: GetPaymentWithLeaseDetails :one
SELECT 
    p.*,
    l.monthly_rent_value,
    l.payment_due_day,
    u.number as unit_number,
    t.full_name as tenant_name,
    t.phone as tenant_phone
FROM payments p
INNER JOIN leases l ON p.lease_id = l.id
INNER JOIN units u ON l.unit_id = u.id
INNER JOIN tenants t ON l.tenant_id = t.id
WHERE p.id = $1
LIMIT 1;

-- name: ListPaymentsWithLeaseDetails :many
SELECT 
    p.*,
    l.monthly_rent_value,
    l.payment_due_day,
    u.number as unit_number,
    t.full_name as tenant_name,
    t.phone as tenant_phone
FROM payments p
INNER JOIN leases l ON p.lease_id = l.id
INNER JOIN units u ON l.unit_id = u.id
INNER JOIN tenants t ON l.tenant_id = t.id
ORDER BY p.due_date DESC;

-- name: GetTotalPaidByLease :one
SELECT COALESCE(SUM(amount), 0) as total
FROM payments
WHERE lease_id = $1
  AND status = 'paid';

-- name: GetPendingAmountByLease :one
SELECT COALESCE(SUM(amount), 0) as total
FROM payments
WHERE lease_id = $1
  AND status IN ('pending', 'overdue');