// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: get_images_for_author_did.sql

package search_queries

import (
	"context"
)

const getImagesForAuthorDID = `-- name: GetImagesForAuthorDID :many
SELECT cid,
    post_id,
    author_did,
    alt_text,
    mime_type,
    created_at,
    cv_completed,
    cv_run_at,
    cv_classes
FROM images
WHERE author_did = $1
`

func (q *Queries) GetImagesForAuthorDID(ctx context.Context, authorDid string) ([]Image, error) {
	rows, err := q.query(ctx, q.getImagesForAuthorDIDStmt, getImagesForAuthorDID, authorDid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Image
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.Cid,
			&i.PostID,
			&i.AuthorDid,
			&i.AltText,
			&i.MimeType,
			&i.CreatedAt,
			&i.CvCompleted,
			&i.CvRunAt,
			&i.CvClasses,
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
