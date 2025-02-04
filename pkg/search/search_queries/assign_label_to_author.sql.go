// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: assign_label_to_author.sql

package search_queries

import (
	"context"
)

const assignLabelToAuthor = `-- name: AssignLabelToAuthor :exec
INSERT INTO author_labels (author_did, label_id)
VALUES ($1, $2::bigint)
ON CONFLICT (author_did, label_id) DO NOTHING
`

type AssignLabelToAuthorParams struct {
	AuthorDid string `json:"author_did"`
	LabelID   int64  `json:"label_id"`
}

func (q *Queries) AssignLabelToAuthor(ctx context.Context, arg AssignLabelToAuthorParams) error {
	_, err := q.exec(ctx, q.assignLabelToAuthorStmt, assignLabelToAuthor, arg.AuthorDid, arg.LabelID)
	return err
}
