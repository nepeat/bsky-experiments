// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: get_unindexed_page.sql

package search_queries

import (
	"context"
)

const getUnindexedPostPage = `-- name: GetUnindexedPostPage :many
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
    p.indexed_at
FROM posts p
WHERE p.indexed_at IS NULL
GROUP BY p.id
ORDER BY p.id
LIMIT $1 OFFSET $2
`

type GetUnindexedPostPageParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetUnindexedPostPage(ctx context.Context, arg GetUnindexedPostPageParams) ([]Post, error) {
	rows, err := q.query(ctx, q.getUnindexedPostPageStmt, getUnindexedPostPage, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
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
