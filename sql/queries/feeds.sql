-- name: FeedsCreate :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
--

-- name: FeedsGetList :many
SELECT * 
FROM feeds;
--

-- name: FeedsGetByURL :one
SELECT * 
FROM feeds
WHERE url = $1;
--

-- name: FeedsMarkedFetched :one
UPDATE feeds
SET updated_at = NOW(), last_fetched_at = NOW()
WHERE id = $1
RETURNING *;
--

-- name: FeedsGetNextToFetch :one
SELECT * 
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
--

-- name: FeedsReset :exec
DELETE FROM feeds;
--
