// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: add_labels_to_posts.sql

package search_queries

import (
	"context"
	"encoding/json"
)

const addLabelsToPosts = `-- name: AddLabelsToPosts :exec
INSERT INTO post_labels (post_id, author_did, label)
SELECT post_id,
    author_did,
    x.label
FROM jsonb_to_recordset($1::jsonb) AS x(post_id text, author_did text, label text) ON CONFLICT (post_id, label) DO NOTHING
`

func (q *Queries) AddLabelsToPosts(ctx context.Context, dollar_1 json.RawMessage) error {
	_, err := q.exec(ctx, q.addLabelsToPostsStmt, addLabelsToPosts, dollar_1)
	return err
}
