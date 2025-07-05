-- name: FeedCreate :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: FeedGetList :many
SELECT * FROM feeds;

-- name: FeedReset :exec
DELETE FROM feeds;
