// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: get_posts_from_cluster_view.sql

package search_queries

import (
	"context"
	"database/sql"
	"time"
)

const getPostsPageByClusterAliasFromView = `-- name: GetPostsPageByClusterAliasFromView :many
SELECT DISTINCT ON (h.id) h.id,
    h.text,
    h.parent_post_id,
    h.root_post_id,
    h.author_did,
    h.created_at,
    h.has_embedded_media,
    h.parent_relationship,
    h.sentiment,
    h.sentiment_confidence,
    h.hotness::float as hotness
FROM post_hotness h
WHERE h.cluster_label = $1
    AND h.created_at < $2
ORDER BY h.created_at DESC
LIMIT $3
`

type GetPostsPageByClusterAliasFromViewParams struct {
	ClusterLabel sql.NullString `json:"cluster_label"`
	CreatedAt    time.Time      `json:"created_at"`
	Limit        int32          `json:"limit"`
}

type GetPostsPageByClusterAliasFromViewRow struct {
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

func (q *Queries) GetPostsPageByClusterAliasFromView(ctx context.Context, arg GetPostsPageByClusterAliasFromViewParams) ([]GetPostsPageByClusterAliasFromViewRow, error) {
	rows, err := q.query(ctx, q.getPostsPageByClusterAliasFromViewStmt, getPostsPageByClusterAliasFromView, arg.ClusterLabel, arg.CreatedAt, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsPageByClusterAliasFromViewRow
	for rows.Next() {
		var i GetPostsPageByClusterAliasFromViewRow
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
