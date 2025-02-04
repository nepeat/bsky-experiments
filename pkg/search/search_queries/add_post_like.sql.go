// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: add_post_like.sql

package search_queries

import (
	"context"
	"database/sql"
)

const addLikeToPost = `-- name: AddLikeToPost :exec
INSERT INTO post_likes (post_id, author_did, like_count)
VALUES ($1, $2, 1) ON CONFLICT (post_id) DO
UPDATE
SET like_count = post_likes.like_count + 1
WHERE post_likes.post_id = $1
`

type AddLikeToPostParams struct {
	PostID    string         `json:"post_id"`
	AuthorDid sql.NullString `json:"author_did"`
}

func (q *Queries) AddLikeToPost(ctx context.Context, arg AddLikeToPostParams) error {
	_, err := q.exec(ctx, q.addLikeToPostStmt, addLikeToPost, arg.PostID, arg.AuthorDid)
	return err
}
