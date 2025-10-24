-- name: CreateLeaseRentAdjustment :one
INSERT INTO lease_rent_adjustments (
    id,
    lease_id,
    previous_rent_value,
    new_rent_value,
    adjustment_percentage,
    applied_at,
    reason,
    applied_by,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetLeaseRentAdjustmentByID :one
SELECT * FROM lease_rent_adjustments
WHERE id = $1
LIMIT 1;

-- name: ListLeaseRentAdjustmentsByLeaseID :many
SELECT * FROM lease_rent_adjustments
WHERE lease_id = $1
ORDER BY applied_at DESC;

-- name: GetLatestAdjustmentByLeaseID :one
SELECT * FROM lease_rent_adjustments
WHERE lease_id = $1
ORDER BY applied_at DESC
LIMIT 1;

-- name: CountAdjustmentsByLeaseID :one
SELECT COUNT(*) FROM lease_rent_adjustments
WHERE lease_id = $1;

-- name: DeleteLeaseRentAdjustment :exec
DELETE FROM lease_rent_adjustments
WHERE id = $1;
