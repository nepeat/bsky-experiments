// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: get_posts_with_any_labels_sorted_by_hotness.sql

package search_queries

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const getPostsPageWithAnyPostLabelSortedByHotness = `-- name: GetPostsPageWithAnyPostLabelSortedByHotness :many
SELECT h.id,
       h.text,
       h.parent_post_id,
       h.root_post_id,
       h.author_did,
       h.created_at,
       h.has_embedded_media,
       h.parent_relationship,
       h.sentiment,
       h.sentiment_confidence,
       MAX(h.hotness)::float as hotness
FROM post_hotness h
WHERE h.post_labels && $1::text []
GROUP BY h.id,
       h.text,
       h.parent_post_id,
       h.root_post_id,
       h.author_did,
       h.created_at,
       h.has_embedded_media,
       h.parent_relationship,
       h.sentiment,
       h.sentiment_confidence,
       hotness
HAVING (
              CASE
                     WHEN $2::float = -1 THEN TRUE
                     ELSE MAX(hotness) < $2::float
              END
       )
ORDER BY MAX(hotness) DESC,
       h.id DESC
LIMIT $3
`

type GetPostsPageWithAnyPostLabelSortedByHotnessParams struct {
	Labels []string `json:"labels"`
	Cursor float64  `json:"cursor"`
	Limit  int32    `json:"limit"`
}

type GetPostsPageWithAnyPostLabelSortedByHotnessRow struct {
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
	Hotness             float64         `json:"hotness"`
}

func (q *Queries) GetPostsPageWithAnyPostLabelSortedByHotness(ctx context.Context, arg GetPostsPageWithAnyPostLabelSortedByHotnessParams) ([]GetPostsPageWithAnyPostLabelSortedByHotnessRow, error) {
	rows, err := q.query(ctx, q.getPostsPageWithAnyPostLabelSortedByHotnessStmt, getPostsPageWithAnyPostLabelSortedByHotness, pq.Array(arg.Labels), arg.Cursor, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsPageWithAnyPostLabelSortedByHotnessRow
	for rows.Next() {
		var i GetPostsPageWithAnyPostLabelSortedByHotnessRow
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
			&i.Hotness,
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
