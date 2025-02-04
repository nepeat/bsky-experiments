// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: add_post_label.sql

package search_queries

import (
	"context"
)

const addPostLabel = `-- name: AddPostLabel :exec
INSERT INTO post_labels (post_id, author_did, label)
VALUES ($1, $2, $3) ON CONFLICT (post_id, label) DO NOTHING
`

type AddPostLabelParams struct {
	PostID    string `json:"post_id"`
	AuthorDid string `json:"author_did"`
	Label     string `json:"label"`
}

func (q *Queries) AddPostLabel(ctx context.Context, arg AddPostLabelParams) error {
	_, err := q.exec(ctx, q.addPostLabelStmt, addPostLabel, arg.PostID, arg.AuthorDid, arg.Label)
	return err
}
