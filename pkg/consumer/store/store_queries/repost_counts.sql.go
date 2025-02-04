// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: repost_counts.sql

package store_queries

import (
	"context"
)

const decrementRepostCountByN = `-- name: DecrementRepostCountByN :exec
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
INSERT INTO repost_counts (subject_id, num_reposts)
SELECT id,
    - $3::int
FROM subj ON CONFLICT (subject_id) DO
UPDATE
SET num_reposts = GREATEST(
        0,
        repost_counts.num_reposts + EXCLUDED.num_reposts
    )
`

type DecrementRepostCountByNParams struct {
	ActorDid   string `json:"actor_did"`
	Rkey       string `json:"rkey"`
	NumReposts int32  `json:"num_reposts"`
	Collection string `json:"collection"`
}

func (q *Queries) DecrementRepostCountByN(ctx context.Context, arg DecrementRepostCountByNParams) error {
	_, err := q.exec(ctx, q.decrementRepostCountByNStmt, decrementRepostCountByN,
		arg.ActorDid,
		arg.Rkey,
		arg.NumReposts,
		arg.Collection,
	)
	return err
}

const deleteRepostCount = `-- name: DeleteRepostCount :exec
DELETE FROM repost_counts
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

type DeleteRepostCountParams struct {
	ActorDid   string `json:"actor_did"`
	Rkey       string `json:"rkey"`
	Collection string `json:"collection"`
}

func (q *Queries) DeleteRepostCount(ctx context.Context, arg DeleteRepostCountParams) error {
	_, err := q.exec(ctx, q.deleteRepostCountStmt, deleteRepostCount, arg.ActorDid, arg.Rkey, arg.Collection)
	return err
}

const getRepostCount = `-- name: GetRepostCount :one
SELECT nlc.subject_id, nlc.num_reposts, nlc.updated_at
FROM repost_counts nlc
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

type GetRepostCountParams struct {
	ActorDid   string `json:"actor_did"`
	Rkey       string `json:"rkey"`
	Collection string `json:"collection"`
}

func (q *Queries) GetRepostCount(ctx context.Context, arg GetRepostCountParams) (RepostCount, error) {
	row := q.queryRow(ctx, q.getRepostCountStmt, getRepostCount, arg.ActorDid, arg.Rkey, arg.Collection)
	var i RepostCount
	err := row.Scan(&i.SubjectID, &i.NumReposts, &i.UpdatedAt)
	return i, err
}

const incrementRepostCountByN = `-- name: IncrementRepostCountByN :exec
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
INSERT INTO repost_counts (subject_id, num_reposts)
SELECT id,
    $3
FROM subj ON CONFLICT (subject_id) DO
UPDATE
SET num_reposts = repost_counts.num_reposts + EXCLUDED.num_reposts
`

type IncrementRepostCountByNParams struct {
	ActorDid   string `json:"actor_did"`
	Rkey       string `json:"rkey"`
	NumReposts int64  `json:"num_reposts"`
	Collection string `json:"collection"`
}

func (q *Queries) IncrementRepostCountByN(ctx context.Context, arg IncrementRepostCountByNParams) error {
	_, err := q.exec(ctx, q.incrementRepostCountByNStmt, incrementRepostCountByN,
		arg.ActorDid,
		arg.Rkey,
		arg.NumReposts,
		arg.Collection,
	)
	return err
}
