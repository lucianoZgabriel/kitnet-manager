-- name: CreateLease :one
INSERT INTO leases (
    id,
    unit_id,
    tenant_id,
    contract_signed_date,
    start_date,
    end_date,
    payment_due_day,
    monthly_rent_value,
    painting_fee_total,
    painting_fee_installments,
    painting_fee_paid,
    status,
    parent_lease_id,
    generation,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
) RETURNING *;

-- name: GetLeaseByID :one
SELECT * FROM leases
WHERE id = $1
LIMIT 1;

-- name: ListLeases :many
SELECT * FROM leases
ORDER BY created_at DESC;

-- name: ListLeasesByStatus :many
SELECT * FROM leases
WHERE status = $1
ORDER BY created_at DESC;

-- name: ListLeasesByUnitID :many
SELECT * FROM leases
WHERE unit_id = $1
ORDER BY created_at DESC;

-- name: ListLeasesByTenantID :many
SELECT * FROM leases
WHERE tenant_id = $1
ORDER BY created_at DESC;

-- name: GetActiveLeaseByUnitID :one
SELECT * FROM leases
WHERE unit_id = $1 AND status = 'active'
LIMIT 1;

-- name: GetActiveLeaseByTenantID :one
SELECT * FROM leases
WHERE tenant_id = $1 AND status = 'active'
LIMIT 1;

-- name: GetExpiringSoonLeases :many
SELECT * FROM leases
WHERE status = 'active' 
  AND end_date <= CURRENT_DATE + INTERVAL '45 days'
  AND end_date > CURRENT_DATE
ORDER BY end_date ASC;

-- name: UpdateLease :one
UPDATE leases
SET
    unit_id = $2,
    tenant_id = $3,
    contract_signed_date = $4,
    start_date = $5,
    end_date = $6,
    payment_due_day = $7,
    monthly_rent_value = $8,
    painting_fee_total = $9,
    painting_fee_installments = $10,
    painting_fee_paid = $11,
    status = $12,
    parent_lease_id = $13,
    generation = $14,
    updated_at = $15
WHERE id = $1
RETURNING *;

-- name: UpdateLeaseStatus :one
UPDATE leases
SET
    status = $2,
    updated_at = $3
WHERE id = $1
RETURNING *;

-- name: UpdatePaintingFeePaid :one
UPDATE leases
SET
    painting_fee_paid = $2,
    updated_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteLease :exec
DELETE FROM leases
WHERE id = $1;

-- name: CountLeases :one
SELECT COUNT(*) FROM leases;

-- name: CountLeasesByStatus :one
SELECT COUNT(*) FROM leases
WHERE status = $1;

-- name: GetLeaseWithDetails :one
SELECT 
    l.*,
    u.number as unit_number,
    u.floor as unit_floor,
    t.full_name as tenant_name,
    t.cpf as tenant_cpf,
    t.phone as tenant_phone
FROM leases l
INNER JOIN units u ON l.unit_id = u.id
INNER JOIN tenants t ON l.tenant_id = t.id
WHERE l.id = $1
LIMIT 1;

-- name: ListLeasesWithDetails :many
SELECT 
    l.*,
    u.number as unit_number,
    u.floor as unit_floor,
    t.full_name as tenant_name,
    t.cpf as tenant_cpf,
    t.phone as tenant_phone
FROM leases l
INNER JOIN units u ON l.unit_id = u.id
INNER JOIN tenants t ON l.tenant_id = t.id
ORDER BY l.created_at DESC;