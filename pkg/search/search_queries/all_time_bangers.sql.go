// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: all_time_bangers.sql

package search_queries

import (
	"context"
	"database/sql"
	"time"
)

const getAllTimeBangers = `-- name: GetAllTimeBangers :many
SELECT p.id,
    p.text,
    p.parent_post_id,
    p.root_post_id,
    p.author_did,
    p.created_at,
    p.has_embedded_media,
    p.parent_relationship,
    p.sentiment,
    p.sentiment_confidence,
    p.indexed_at,
    l.like_count
FROM posts p
    JOIN post_likes l ON l.post_id = p.id
ORDER BY l.like_count DESC
LIMIT $2 OFFSET $1
`

type GetAllTimeBangersParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

type GetAllTimeBangersRow struct {
	ID                  string          `json:"id"`
	Text                string          `json:"text"`
	ParentPostID        sql.NullString  `json:"parent_post_id"`
	RootPostID          sql.NullString  `json:"root_post_id"`
	AuthorDid           string          `json:"author_did"`
	CreatedAt           time.Time       `json:"created_at"`
	HasEmbeddedMedia    bool            `json:"has_embedded_media"`
	ParentRelationship  sql.NullString  `json:"parent_relationship"`
	Sentiment           sql.NullString  `json:"sentiment"`
	SentimentConfidence sql.NullFloat64 `json:"sentiment_confidence"`
	IndexedAt           sql.NullTime    `json:"indexed_at"`
	LikeCount           int64           `json:"like_count"`
}

func (q *Queries) GetAllTimeBangers(ctx context.Context, arg GetAllTimeBangersParams) ([]GetAllTimeBangersRow, error) {
	rows, err := q.query(ctx, q.getAllTimeBangersStmt, getAllTimeBangers, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllTimeBangersRow
	for rows.Next() {
		var i GetAllTimeBangersRow
		if err := rows.Scan(
			&i.ID,
			&i.Text,
			&i.ParentPostID,
			&i.RootPostID,
			&i.AuthorDid,
			&i.CreatedAt,
			&i.HasEmbeddedMedia,
			&i.ParentRelationship,
			&i.Sentiment,
			&i.SentimentConfidence,
			&i.IndexedAt,
			&i.LikeCount,
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
