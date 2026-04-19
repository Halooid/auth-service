-- name: GetUserByExternalID :one
SELECT * FROM users
WHERE external_id = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    external_id, email, username, tenant_id, first_name, last_name
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET 
    email = $2,
    username = $3,
    tenant_id = $4,
    first_name = $5,
    last_name = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
