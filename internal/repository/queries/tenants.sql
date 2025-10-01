-- name: CreateTenant :one
INSERT INTO tenants (
    id,
    full_name,
    cpf,
    phone,
    email,
    id_document_type,
    id_document_number,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetTenantByID :one
SELECT * FROM tenants
WHERE id = $1
LIMIT 1;

-- name: GetTenantByCPF :one
SELECT * FROM tenants
WHERE cpf = $1
LIMIT 1;

-- name: ListTenants :many
SELECT * FROM tenants
ORDER BY full_name ASC;

-- name: SearchTenantsByName :many
SELECT * FROM tenants
WHERE full_name ILIKE '%' || $1 || '%'
ORDER BY full_name ASC;

-- name: UpdateTenant :one
UPDATE tenants
SET
    full_name = $2,
    phone = $3,
    email = $4,
    id_document_type = $5,
    id_document_number = $6,
    updated_at = $7
WHERE id = $1
RETURNING *;

-- name: DeleteTenant :exec
DELETE FROM tenants
WHERE id = $1;

-- name: CountTenants :one
SELECT COUNT(*) FROM tenants;

-- name: TenantExistsByCPF :one
SELECT EXISTS(
    SELECT 1 FROM tenants WHERE cpf = $1
) AS exists;