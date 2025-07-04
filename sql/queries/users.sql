-- name: UserCreate :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: UserGetByName :one
SELECT * FROM users WHERE name = $1;

-- name: UserGetList :many
SELECT * FROM users;

-- name: UserReset :exec
DELETE FROM users;
