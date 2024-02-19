// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: reposts.sql

package store_queries

import (
	"context"
	"database/sql"
	"time"
)

const createRepost = `-- name: CreateRepost :exec
WITH collection_ins AS (
    INSERT INTO collections (name)
    VALUES ($4) ON CONFLICT (name) DO NOTHING
    RETURNING id
),
subject_ins AS (
    INSERT INTO subjects (actor_did, rkey, col)
    VALUES (
            $5,
            $6,
            COALESCE(
                (
                    SELECT id
                    FROM collection_ins
                ),
                (
                    SELECT id
                    FROM collections
                    WHERE name = $4
                )
            )
        ) ON CONFLICT (actor_did, col, rkey) DO
    UPDATE
    SET actor_did = EXCLUDED.actor_did
    RETURNING id
)
INSERT INTO reposts (actor_did, rkey, subj, created_at)
SELECT $1,
    $2,
    subject_ins.id,
    $3
FROM subject_ins
`

type CreateRepostParams struct {
	ActorDid        string       `json:"actor_did"`
	Rkey            string       `json:"rkey"`
	CreatedAt       sql.NullTime `json:"created_at"`
	Collection      string       `json:"collection"`
	SubjectActorDid string       `json:"subject_actor_did"`
	SubjectRkey     string       `json:"subject_rkey"`
}

func (q *Queries) CreateRepost(ctx context.Context, arg CreateRepostParams) error {
	_, err := q.exec(ctx, q.createRepostStmt, createRepost,
		arg.ActorDid,
		arg.Rkey,
		arg.CreatedAt,
		arg.Collection,
		arg.SubjectActorDid,
		arg.SubjectRkey,
	)
	return err
}

const deleteRepost = `-- name: DeleteRepost :exec
DELETE FROM reposts
WHERE actor_did = $1
    AND rkey = $2
`

type DeleteRepostParams struct {
	ActorDid string `json:"actor_did"`
	Rkey     string `json:"rkey"`
}

func (q *Queries) DeleteRepost(ctx context.Context, arg DeleteRepostParams) error {
	_, err := q.exec(ctx, q.deleteRepostStmt, deleteRepost, arg.ActorDid, arg.Rkey)
	return err
}

const getRepost = `-- name: GetRepost :one
SELECT l.actor_did, l.rkey, l.subj, l.created_at, l.inserted_at,
    s.actor_did AS subject_actor_did,
    c.name AS subject_namespace,
    s.rkey AS subject_rkey
FROM reposts l
    JOIN subjects s ON l.subj = s.id
    JOIN collections c ON s.col = c.id
WHERE l.actor_did = $1
    AND l.rkey = $2
LIMIT 1
`

type GetRepostParams struct {
	ActorDid string `json:"actor_did"`
	Rkey     string `json:"rkey"`
}

type GetRepostRow struct {
	ActorDid         string       `json:"actor_did"`
	Rkey             string       `json:"rkey"`
	Subj             int64        `json:"subj"`
	CreatedAt        sql.NullTime `json:"created_at"`
	InsertedAt       time.Time    `json:"inserted_at"`
	SubjectActorDid  string       `json:"subject_actor_did"`
	SubjectNamespace string       `json:"subject_namespace"`
	SubjectRkey      string       `json:"subject_rkey"`
}

func (q *Queries) GetRepost(ctx context.Context, arg GetRepostParams) (GetRepostRow, error) {
	row := q.queryRow(ctx, q.getRepostStmt, getRepost, arg.ActorDid, arg.Rkey)
	var i GetRepostRow
	err := row.Scan(
		&i.ActorDid,
		&i.Rkey,
		&i.Subj,
		&i.CreatedAt,
		&i.InsertedAt,
		&i.SubjectActorDid,
		&i.SubjectNamespace,
		&i.SubjectRkey,
	)
	return i, err
}

const getRepostsByActor = `-- name: GetRepostsByActor :many
SELECT l.actor_did, l.rkey, l.subj, l.created_at, l.inserted_at,
    s.actor_did AS subject_actor_did,
    c.name AS subject_namespace,
    s.rkey AS subject_rkey
FROM reposts l
    JOIN subjects s ON l.subj = s.id
    JOIN collections c ON s.col = c.id
WHERE l.actor_did = $1
ORDER BY l.created_at DESC
LIMIT $2 OFFSET $3
`

type GetRepostsByActorParams struct {
	ActorDid string `json:"actor_did"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

type GetRepostsByActorRow struct {
	ActorDid         string       `json:"actor_did"`
	Rkey             string       `json:"rkey"`
	Subj             int64        `json:"subj"`
	CreatedAt        sql.NullTime `json:"created_at"`
	InsertedAt       time.Time    `json:"inserted_at"`
	SubjectActorDid  string       `json:"subject_actor_did"`
	SubjectNamespace string       `json:"subject_namespace"`
	SubjectRkey      string       `json:"subject_rkey"`
}

func (q *Queries) GetRepostsByActor(ctx context.Context, arg GetRepostsByActorParams) ([]GetRepostsByActorRow, error) {
	rows, err := q.query(ctx, q.getRepostsByActorStmt, getRepostsByActor, arg.ActorDid, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRepostsByActorRow
	for rows.Next() {
		var i GetRepostsByActorRow
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.Subj,
			&i.CreatedAt,
			&i.InsertedAt,
			&i.SubjectActorDid,
			&i.SubjectNamespace,
			&i.SubjectRkey,
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

const getRepostsBySubject = `-- name: GetRepostsBySubject :many
SELECT l.actor_did, l.rkey, l.subj, l.created_at, l.inserted_at,
    s.actor_did AS subject_actor_did,
    c.name AS subject_namespace,
    s.rkey AS subject_rkey
FROM reposts l
    JOIN subjects s ON l.subj = s.id
    JOIN collections c ON s.col = c.id
WHERE s.actor_did = $1
    AND c.name = $2
    AND s.rkey = $3
ORDER BY l.created_at DESC
LIMIT $4 OFFSET $5
`

type GetRepostsBySubjectParams struct {
	ActorDid string `json:"actor_did"`
	Name     string `json:"name"`
	Rkey     string `json:"rkey"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

type GetRepostsBySubjectRow struct {
	ActorDid         string       `json:"actor_did"`
	Rkey             string       `json:"rkey"`
	Subj             int64        `json:"subj"`
	CreatedAt        sql.NullTime `json:"created_at"`
	InsertedAt       time.Time    `json:"inserted_at"`
	SubjectActorDid  string       `json:"subject_actor_did"`
	SubjectNamespace string       `json:"subject_namespace"`
	SubjectRkey      string       `json:"subject_rkey"`
}

func (q *Queries) GetRepostsBySubject(ctx context.Context, arg GetRepostsBySubjectParams) ([]GetRepostsBySubjectRow, error) {
	rows, err := q.query(ctx, q.getRepostsBySubjectStmt, getRepostsBySubject,
		arg.ActorDid,
		arg.Name,
		arg.Rkey,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRepostsBySubjectRow
	for rows.Next() {
		var i GetRepostsBySubjectRow
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.Subj,
			&i.CreatedAt,
			&i.InsertedAt,
			&i.SubjectActorDid,
			&i.SubjectNamespace,
			&i.SubjectRkey,
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
