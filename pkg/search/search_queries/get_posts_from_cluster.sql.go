// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: get_posts_from_cluster.sql

package search_queries

import (
	"context"
	"database/sql"
	"time"
)

const getPostsPageByClusterAlias = `-- name: GetPostsPageByClusterAlias :many
SELECT p.id,
      p.text,
      p.parent_post_id,
      p.root_post_id,
      p.author_did,
      p.created_at,
      p.has_embedded_media,
      p.parent_relationship,
      p.sentiment,
      p.sentiment_confidence
FROM posts p
      JOIN author_clusters ON p.author_did = author_clusters.author_did
      JOIN clusters ON author_clusters.cluster_id = clusters.id
WHERE clusters.lookup_alias = $1
      AND p.created_at < $2
      AND p.created_at >= NOW() - make_interval(hours := CAST($4 AS INT))
ORDER BY p.created_at DESC
LIMIT $3
`

type GetPostsPageByClusterAliasParams struct {
	LookupAlias string    `json:"lookup_alias"`
	CreatedAt   time.Time `json:"created_at"`
	Limit       int32     `json:"limit"`
	HoursAgo    int32     `json:"hours_ago"`
}

type GetPostsPageByClusterAliasRow struct {
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
}

func (q *Queries) GetPostsPageByClusterAlias(ctx context.Context, arg GetPostsPageByClusterAliasParams) ([]GetPostsPageByClusterAliasRow, error) {
	rows, err := q.query(ctx, q.getPostsPageByClusterAliasStmt, getPostsPageByClusterAlias,
		arg.LookupAlias,
		arg.CreatedAt,
		arg.Limit,
		arg.HoursAgo,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsPageByClusterAliasRow
	for rows.Next() {
		var i GetPostsPageByClusterAliasRow
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
