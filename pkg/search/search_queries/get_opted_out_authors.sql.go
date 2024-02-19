// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: get_opted_out_authors.sql

package search_queries

import (
	"context"
)

const getOptedOutAuthors = `-- name: GetOptedOutAuthors :many
SELECT did,
    handle,
    cluster_opt_out
FROM authors
WHERE cluster_opt_out = TRUE
`

func (q *Queries) GetOptedOutAuthors(ctx context.Context) ([]Author, error) {
	rows, err := q.query(ctx, q.getOptedOutAuthorsStmt, getOptedOutAuthors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Author
	for rows.Next() {
		var i Author
		if err := rows.Scan(&i.Did, &i.Handle, &i.ClusterOptOut); err != nil {
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
