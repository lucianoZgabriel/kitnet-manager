-- name: CreateUnit :one
INSERT INTO units (
    id,
    number,
    floor,
    status,
    is_renovated,
    base_rent_value,
    renovated_rent_value,
    current_rent_value,
    notes,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetUnitByID :one
SELECT * FROM units
WHERE id = $1
LIMIT 1;

-- name: GetUnitByNumber :one
SELECT * FROM units
WHERE number = $1
LIMIT 1;

-- name: ListUnits :many
SELECT * FROM units
ORDER BY floor ASC, number ASC;

-- name: ListUnitsByStatus :many
SELECT * FROM units
WHERE status = $1
ORDER BY floor ASC, number ASC;

-- name: ListUnitsByFloor :many
SELECT * FROM units
WHERE floor = $1
ORDER BY number ASC;

-- name: ListAvailableUnits :many
SELECT * FROM units
WHERE status = 'available'
ORDER BY floor ASC, number ASC;

-- name: UpdateUnit :one
UPDATE units
SET
    number = $2,
    floor = $3,
    status = $4,
    is_renovated = $5,
    base_rent_value = $6,
    renovated_rent_value = $7,
    current_rent_value = $8,
    notes = $9,
    updated_at = $10
WHERE id = $1
RETURNING *;

-- name: UpdateUnitStatus :one
UPDATE units
SET
    status = $2,
    updated_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUnit :exec
DELETE FROM units
WHERE id = $1;

-- name: CountUnits :one
SELECT COUNT(*) FROM units;

-- name: CountUnitsByStatus :one
SELECT COUNT(*) FROM units
WHERE status = $1;