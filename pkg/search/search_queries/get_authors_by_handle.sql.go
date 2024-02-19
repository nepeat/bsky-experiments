// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: get_authors_by_handle.sql

package search_queries

import (
	"context"
)

const getAuthorsByHandle = `-- name: GetAuthorsByHandle :many
SELECT did, handle
FROM authors
WHERE handle = $1
`

type GetAuthorsByHandleRow struct {
	Did    string `json:"did"`
	Handle string `json:"handle"`
}

func (q *Queries) GetAuthorsByHandle(ctx context.Context, handle string) ([]GetAuthorsByHandleRow, error) {
	rows, err := q.query(ctx, q.getAuthorsByHandleStmt, getAuthorsByHandle, handle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAuthorsByHandleRow
	for rows.Next() {
		var i GetAuthorsByHandleRow
		if err := rows.Scan(&i.Did, &i.Handle); err != nil {
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
