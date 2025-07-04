-- name: FeedCreate :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: FeedGetList :many
SELECT feeds.name, feeds.url, users.name as user_name
FROM feeds
LEFT JOIN users
ON feeds.user_id = users.id
ORDER BY user_name;


-- name: FeedReset :exec
DELETE FROM feeds;
