// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: like_counts.sql

package store_queries

import (
	"context"
	"database/sql"
	"time"
)

const createLikeCount = `-- name: CreateLikeCount :exec
INSERT INTO like_counts (
        subject_id,
        num_likes,
        updated_at,
        subject_created_at
    )
VALUES ($1, $2, $3, $4)
`

type CreateLikeCountParams struct {
	SubjectID        int64        `json:"subject_id"`
	NumLikes         int64        `json:"num_likes"`
	UpdatedAt        time.Time    `json:"updated_at"`
	SubjectCreatedAt sql.NullTime `json:"subject_created_at"`
}

func (q *Queries) CreateLikeCount(ctx context.Context, arg CreateLikeCountParams) error {
	_, err := q.exec(ctx, q.createLikeCountStmt, createLikeCount,
		arg.SubjectID,
		arg.NumLikes,
		arg.UpdatedAt,
		arg.SubjectCreatedAt,
	)
	return err
}

const decrementLikeCountByN = `-- name: DecrementLikeCountByN :exec
WITH subj AS (
    SELECT id
    FROM subjects
    WHERE actor_did = $1
        AND col = (
            SELECT id
            FROM collections
            WHERE name = $4
        )
        AND rkey = $2
)
INSERT INTO like_counts (subject_id, num_likes)
SELECT id,
    - $3::int
FROM subj ON CONFLICT (subject_id) DO
UPDATE
SET num_likes = GREATEST(
        0,
        like_counts.num_likes + EXCLUDED.num_likes
    )
`

type DecrementLikeCountByNParams struct {
	ActorDid   string `json:"actor_did"`
	Rkey       string `json:"rkey"`
	NumLikes   int32  `json:"num_likes"`
	Collection string `json:"collection"`
}

func (q *Queries) DecrementLikeCountByN(ctx context.Context, arg DecrementLikeCountByNParams) error {
	_, err := q.exec(ctx, q.decrementLikeCountByNStmt, decrementLikeCountByN,
		arg.ActorDid,
		arg.Rkey,
		arg.NumLikes,
		arg.Collection,
	)
	return err
}

const deleteLikeCount = `-- name: DeleteLikeCount :exec
DELETE FROM like_counts
WHERE subject_id IN (
        SELECT id
        FROM subjects
        WHERE actor_did = $1
            AND col = (
                SELECT id
                FROM collections
                WHERE name = $3
            )
            AND rkey = $2
    )
`

type DeleteLikeCountParams struct {
	ActorDid   string `json:"actor_did"`
	Rkey       string `json:"rkey"`
	Collection string `json:"collection"`
}

func (q *Queries) DeleteLikeCount(ctx context.Context, arg DeleteLikeCountParams) error {
	_, err := q.exec(ctx, q.deleteLikeCountStmt, deleteLikeCount, arg.ActorDid, arg.Rkey, arg.Collection)
	return err
}

const getLikeCount = `-- name: GetLikeCount :one
SELECT nlc.subject_id, nlc.num_likes, nlc.updated_at, nlc.subject_created_at
FROM like_counts nlc
    JOIN subjects s ON nlc.subject_id = s.id
WHERE s.actor_did = $1
    AND s.col = (
        SELECT id
        FROM collections
        WHERE name = $3
    )
    AND s.rkey = $2
LIMIT 1
`

type GetLikeCountParams struct {
	ActorDid   string `json:"actor_did"`
	Rkey       string `json:"rkey"`
	Collection string `json:"collection"`
}

func (q *Queries) GetLikeCount(ctx context.Context, arg GetLikeCountParams) (LikeCount, error) {
	row := q.queryRow(ctx, q.getLikeCountStmt, getLikeCount, arg.ActorDid, arg.Rkey, arg.Collection)
	var i LikeCount
	err := row.Scan(
		&i.SubjectID,
		&i.NumLikes,
		&i.UpdatedAt,
		&i.SubjectCreatedAt,
	)
	return i, err
}

const incrementLikeCountByN = `-- name: IncrementLikeCountByN :exec
WITH subj AS (
    SELECT id
    FROM subjects
    WHERE actor_did = $1
        AND col = (
            SELECT id
            FROM collections
            WHERE name = $4
        )
        AND rkey = $2
)
INSERT INTO like_counts (subject_id, num_likes)
SELECT id,
    $3
FROM subj ON CONFLICT (subject_id) DO
UPDATE
SET num_likes = like_counts.num_likes + EXCLUDED.num_likes
`

type IncrementLikeCountByNParams struct {
	ActorDid   string `json:"actor_did"`
	Rkey       string `json:"rkey"`
	NumLikes   int64  `json:"num_likes"`
	Collection string `json:"collection"`
}

func (q *Queries) IncrementLikeCountByN(ctx context.Context, arg IncrementLikeCountByNParams) error {
	_, err := q.exec(ctx, q.incrementLikeCountByNStmt, incrementLikeCountByN,
		arg.ActorDid,
		arg.Rkey,
		arg.NumLikes,
		arg.Collection,
	)
	return err
}

const incrementLikeCountByNWithSubject = `-- name: IncrementLikeCountByNWithSubject :exec
INSERT INTO like_counts (subject_id, num_likes)
VALUES ($1, $2)
ON CONFLICT (subject_id) DO
UPDATE
SET num_likes = like_counts.num_likes + EXCLUDED.num_likes
`

type IncrementLikeCountByNWithSubjectParams struct {
	SubjectID int64 `json:"subject_id"`
	NumLikes  int64 `json:"num_likes"`
}

func (q *Queries) IncrementLikeCountByNWithSubject(ctx context.Context, arg IncrementLikeCountByNWithSubjectParams) error {
	_, err := q.exec(ctx, q.incrementLikeCountByNWithSubjectStmt, incrementLikeCountByNWithSubject, arg.SubjectID, arg.NumLikes)
	return err
}
