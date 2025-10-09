-- name: CreateUser :one
INSERT INTO users (
    id,
    username,
    password_hash,
    role,
    is_active,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: ListUsersByRole :many
SELECT * FROM users
WHERE role = $1
ORDER BY username ASC;

-- name: UpdateUser :one
UPDATE users
SET
    role = $2,
    is_active = $3,
    updated_at = $4
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET
    password_hash = $2,
    updated_at = $3
WHERE id = $1
RETURNING *;

-- name: UpdateLastLogin :one
UPDATE users
SET
    last_login_at = $2,
    updated_at = $3
WHERE id = $1
RETURNING *;

-- name: DeactivateUser :exec
UPDATE users
SET
    is_active = false,
    updated_at = $2
WHERE id = $1;

-- name: ActivateUser :exec
UPDATE users
SET
    is_active = true,
    updated_at = $2
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UserExistsByUsername :one
SELECT EXISTS(
    SELECT 1 FROM users WHERE username = $1
) AS exists;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CountActiveUsers :one
SELECT COUNT(*) FROM users
WHERE is_active = true;