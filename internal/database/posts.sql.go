// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const postsCreate = `-- name: PostsCreate :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, created_at, updated_at, title, url, description, published_at, feed_id
`

type PostsCreateParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt sql.NullTime
	FeedID      uuid.UUID
}

func (q *Queries) PostsCreate(ctx context.Context, arg PostsCreateParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, postsCreate,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}

const postsGetForUser = `-- name: PostsGetForUser :many

SELECT posts.id, posts.created_at, posts.updated_at, posts.title, posts.url, posts.description, posts.published_at, posts.feed_id, feeds.name AS feed_name
FROM posts
JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
JOIN feeds ON posts.feed_id = feeds.id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2
`

type PostsGetForUserParams struct {
	UserID uuid.UUID
	Limit  int32
}

type PostsGetForUserRow struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt sql.NullTime
	FeedID      uuid.UUID
	FeedName    string
}

func (q *Queries) PostsGetForUser(ctx context.Context, arg PostsGetForUserParams) ([]PostsGetForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, postsGetForUser, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PostsGetForUserRow
	for rows.Next() {
		var i PostsGetForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
			&i.FeedName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
