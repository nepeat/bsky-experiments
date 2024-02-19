// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: get_labels_for_author.sql

package search_queries

import (
	"context"
)

const getLabelsForAuthor = `-- name: GetLabelsForAuthor :many
SELECT labels.id, labels.name, labels.lookup_alias
FROM authors
JOIN author_labels ON authors.did = author_labels.author_did
JOIN labels on author_labels.label_id = labels.id
WHERE authors.did = $1
`

type GetLabelsForAuthorRow struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	LookupAlias string `json:"lookup_alias"`
}

func (q *Queries) GetLabelsForAuthor(ctx context.Context, authorDid string) ([]GetLabelsForAuthorRow, error) {
	rows, err := q.query(ctx, q.getLabelsForAuthorStmt, getLabelsForAuthor, authorDid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLabelsForAuthorRow
	for rows.Next() {
		var i GetLabelsForAuthorRow
		if err := rows.Scan(&i.ID, &i.Name, &i.LookupAlias); err != nil {
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
