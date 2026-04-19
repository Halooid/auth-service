-- name: GetUserByExternalID :one
SELECT * FROM users
WHERE external_id = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    external_id, email, username, tenant_id
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET 
    email = $2,
    username = $3,
    tenant_id = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
