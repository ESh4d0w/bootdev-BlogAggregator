// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: feeds.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const feedsCreate = `-- name: FeedsCreate :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, updated_at, name, url, user_id
`

type FeedsCreateParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	UserID    uuid.UUID
}

func (q *Queries) FeedsCreate(ctx context.Context, arg FeedsCreateParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, feedsCreate,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}

const feedsGetByURL = `-- name: FeedsGetByURL :one
SELECT id, created_at, updated_at, name, url, user_id FROM feeds
WHERE url = $1
`

func (q *Queries) FeedsGetByURL(ctx context.Context, url string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, feedsGetByURL, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}

const feedsGetList = `-- name: FeedsGetList :many
SELECT id, created_at, updated_at, name, url, user_id FROM feeds
`

func (q *Queries) FeedsGetList(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, feedsGetList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
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

const feedsReset = `-- name: FeedsReset :exec
DELETE FROM feeds
`

func (q *Queries) FeedsReset(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, feedsReset)
	return err
}
