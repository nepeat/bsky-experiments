// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: likes.sql

package store_queries

import (
	"context"
	"database/sql"
	"time"
)

const createLike = `-- name: CreateLike :exec
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
INSERT INTO likes (actor_did, rkey, subj, created_at)
SELECT $1,
    $2,
    subject_ins.id,
    $3
FROM subject_ins
`

type CreateLikeParams struct {
	ActorDid        string       `json:"actor_did"`
	Rkey            string       `json:"rkey"`
	CreatedAt       sql.NullTime `json:"created_at"`
	Collection      string       `json:"collection"`
	SubjectActorDid string       `json:"subject_actor_did"`
	SubjectRkey     string       `json:"subject_rkey"`
}

func (q *Queries) CreateLike(ctx context.Context, arg CreateLikeParams) error {
	_, err := q.exec(ctx, q.createLikeStmt, createLike,
		arg.ActorDid,
		arg.Rkey,
		arg.CreatedAt,
		arg.Collection,
		arg.SubjectActorDid,
		arg.SubjectRkey,
	)
	return err
}

const deleteLike = `-- name: DeleteLike :exec
DELETE FROM likes
WHERE actor_did = $1
    AND rkey = $2
`

type DeleteLikeParams struct {
	ActorDid string `json:"actor_did"`
	Rkey     string `json:"rkey"`
}

func (q *Queries) DeleteLike(ctx context.Context, arg DeleteLikeParams) error {
	_, err := q.exec(ctx, q.deleteLikeStmt, deleteLike, arg.ActorDid, arg.Rkey)
	return err
}

const findPotentialFriends = `-- name: FindPotentialFriends :many
WITH user_likes AS (
    SELECT subj
    FROM likes
    WHERE actor_did = $1
        AND created_at > CURRENT_TIMESTAMP - INTERVAL '7 days'
)
SELECT l.actor_did,
    COUNT(l.subj) AS overlap_count
FROM likes l
    JOIN user_likes ul ON l.subj = ul.subj
    LEFT JOIN follows f ON l.actor_did = f.target_did
    AND f.actor_did = $1
    LEFT JOIN blocks b1 ON l.actor_did = b1.target_did
    AND b1.actor_did = $1
    LEFT JOIN blocks b2 ON l.actor_did = b2.actor_did
    AND b2.target_did = $1
WHERE l.actor_did != $1
    AND l.actor_did NOT IN (
        'did:plc:xxno7p4xtpkxtn4ok6prtlcb',
        'did:plc:3tm2l7kcljcgacctmmqru3hj'
    )
    AND l.created_at > CURRENT_TIMESTAMP - INTERVAL '7 days'
    AND f.target_did IS NULL
    AND b1.target_did IS NULL
    AND b2.target_did IS NULL
GROUP BY l.actor_did
ORDER BY overlap_count DESC,
    l.actor_did
LIMIT $2
`

type FindPotentialFriendsParams struct {
	ActorDid string `json:"actor_did"`
	Limit    int32  `json:"limit"`
}

type FindPotentialFriendsRow struct {
	ActorDid     string `json:"actor_did"`
	OverlapCount int64  `json:"overlap_count"`
}

func (q *Queries) FindPotentialFriends(ctx context.Context, arg FindPotentialFriendsParams) ([]FindPotentialFriendsRow, error) {
	rows, err := q.query(ctx, q.findPotentialFriendsStmt, findPotentialFriends, arg.ActorDid, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindPotentialFriendsRow
	for rows.Next() {
		var i FindPotentialFriendsRow
		if err := rows.Scan(&i.ActorDid, &i.OverlapCount); err != nil {
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

const getLike = `-- name: GetLike :one
SELECT l.actor_did, l.rkey, l.subj, l.created_at, l.inserted_at,
    s.actor_did AS subject_actor_did,
    c.name AS subject_namespace,
    s.rkey AS subject_rkey
FROM likes l
    JOIN subjects s ON l.subj = s.id
    JOIN collections c ON s.col = c.id
WHERE l.actor_did = $1
    AND l.rkey = $2
LIMIT 1
`

type GetLikeParams struct {
	ActorDid string `json:"actor_did"`
	Rkey     string `json:"rkey"`
}

type GetLikeRow struct {
	ActorDid         string       `json:"actor_did"`
	Rkey             string       `json:"rkey"`
	Subj             int64        `json:"subj"`
	CreatedAt        sql.NullTime `json:"created_at"`
	InsertedAt       time.Time    `json:"inserted_at"`
	SubjectActorDid  string       `json:"subject_actor_did"`
	SubjectNamespace string       `json:"subject_namespace"`
	SubjectRkey      string       `json:"subject_rkey"`
}

func (q *Queries) GetLike(ctx context.Context, arg GetLikeParams) (GetLikeRow, error) {
	row := q.queryRow(ctx, q.getLikeStmt, getLike, arg.ActorDid, arg.Rkey)
	var i GetLikeRow
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

const getLikesByActor = `-- name: GetLikesByActor :many
SELECT l.actor_did, l.rkey, l.subj, l.created_at, l.inserted_at,
    s.actor_did AS subject_actor_did,
    c.name AS subject_namespace,
    s.rkey AS subject_rkey
FROM likes l
    JOIN subjects s ON l.subj = s.id
    JOIN collections c ON s.col = c.id
WHERE l.actor_did = $1
ORDER BY l.created_at DESC
LIMIT $2 OFFSET $3
`

type GetLikesByActorParams struct {
	ActorDid string `json:"actor_did"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

type GetLikesByActorRow struct {
	ActorDid         string       `json:"actor_did"`
	Rkey             string       `json:"rkey"`
	Subj             int64        `json:"subj"`
	CreatedAt        sql.NullTime `json:"created_at"`
	InsertedAt       time.Time    `json:"inserted_at"`
	SubjectActorDid  string       `json:"subject_actor_did"`
	SubjectNamespace string       `json:"subject_namespace"`
	SubjectRkey      string       `json:"subject_rkey"`
}

func (q *Queries) GetLikesByActor(ctx context.Context, arg GetLikesByActorParams) ([]GetLikesByActorRow, error) {
	rows, err := q.query(ctx, q.getLikesByActorStmt, getLikesByActor, arg.ActorDid, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLikesByActorRow
	for rows.Next() {
		var i GetLikesByActorRow
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

const getLikesBySubject = `-- name: GetLikesBySubject :many
SELECT l.actor_did, l.rkey, l.subj, l.created_at, l.inserted_at,
    s.actor_did AS subject_actor_did,
    c.name AS subject_namespace,
    s.rkey AS subject_rkey
FROM likes l
    JOIN subjects s ON l.subj = s.id
    JOIN collections c ON s.col = c.id
WHERE s.actor_did = $1
    AND c.name = $2
    AND s.rkey = $3
ORDER BY l.created_at DESC
LIMIT $4 OFFSET $5
`

type GetLikesBySubjectParams struct {
	ActorDid string `json:"actor_did"`
	Name     string `json:"name"`
	Rkey     string `json:"rkey"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

type GetLikesBySubjectRow struct {
	ActorDid         string       `json:"actor_did"`
	Rkey             string       `json:"rkey"`
	Subj             int64        `json:"subj"`
	CreatedAt        sql.NullTime `json:"created_at"`
	InsertedAt       time.Time    `json:"inserted_at"`
	SubjectActorDid  string       `json:"subject_actor_did"`
	SubjectNamespace string       `json:"subject_namespace"`
	SubjectRkey      string       `json:"subject_rkey"`
}

func (q *Queries) GetLikesBySubject(ctx context.Context, arg GetLikesBySubjectParams) ([]GetLikesBySubjectRow, error) {
	rows, err := q.query(ctx, q.getLikesBySubjectStmt, getLikesBySubject,
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
	var items []GetLikesBySubjectRow
	for rows.Next() {
		var i GetLikesBySubjectRow
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

const getLikesGivenByActorFromTo = `-- name: GetLikesGivenByActorFromTo :one
WITH p AS MATERIALIZED (
    SELECT p.created_at
    FROM likes l
        JOIN subjects s ON l.subj = s.id
        JOIN posts p ON s.actor_did = p.actor_did
        AND s.rkey = p.rkey
    WHERE l.actor_did = $1
        AND l.created_at > $2
        AND l.created_at < $3
)
SELECT count(*)
FROM p
WHERE p.created_at > $2
    AND p.created_at < $3
`

type GetLikesGivenByActorFromToParams struct {
	ActorDid string       `json:"actor_did"`
	From     sql.NullTime `json:"from"`
	To       sql.NullTime `json:"to"`
}

func (q *Queries) GetLikesGivenByActorFromTo(ctx context.Context, arg GetLikesGivenByActorFromToParams) (int64, error) {
	row := q.queryRow(ctx, q.getLikesGivenByActorFromToStmt, getLikesGivenByActorFromTo, arg.ActorDid, arg.From, arg.To)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getLikesReceivedByActorFromActor = `-- name: GetLikesReceivedByActorFromActor :one
SELECT COUNT(*)
FROM likes
    JOIN subjects ON likes.subj = subjects.id
WHERE subjects.actor_did = $1
    AND likes.actor_did = $2
`

type GetLikesReceivedByActorFromActorParams struct {
	To   string `json:"to"`
	From string `json:"from"`
}

func (q *Queries) GetLikesReceivedByActorFromActor(ctx context.Context, arg GetLikesReceivedByActorFromActorParams) (int64, error) {
	row := q.queryRow(ctx, q.getLikesReceivedByActorFromActorStmt, getLikesReceivedByActorFromActor, arg.To, arg.From)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getTotalLikesGivenByActor = `-- name: GetTotalLikesGivenByActor :one
SELECT COUNT(*)
FROM likes
WHERE actor_did = $1
`

func (q *Queries) GetTotalLikesGivenByActor(ctx context.Context, actorDid string) (int64, error) {
	row := q.queryRow(ctx, q.getTotalLikesGivenByActorStmt, getTotalLikesGivenByActor, actorDid)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getTotalLikesReceivedByActor = `-- name: GetTotalLikesReceivedByActor :one
SELECT SUM(num_likes)
FROM like_counts
    JOIN subjects ON like_counts.subject_id = subjects.id
WHERE subjects.actor_did = $1
`

func (q *Queries) GetTotalLikesReceivedByActor(ctx context.Context, actorDid string) (int64, error) {
	row := q.queryRow(ctx, q.getTotalLikesReceivedByActorStmt, getTotalLikesReceivedByActor, actorDid)
	var sum int64
	err := row.Scan(&sum)
	return sum, err
}

const insertLike = `-- name: InsertLike :exec
INSERT INTO likes (actor_did, rkey, subj, created_at)
VALUES ($1, $2, $3, $4)
`

type InsertLikeParams struct {
	ActorDid  string       `json:"actor_did"`
	Rkey      string       `json:"rkey"`
	Subj      int64        `json:"subj"`
	CreatedAt sql.NullTime `json:"created_at"`
}

func (q *Queries) InsertLike(ctx context.Context, arg InsertLikeParams) error {
	_, err := q.exec(ctx, q.insertLikeStmt, insertLike,
		arg.ActorDid,
		arg.Rkey,
		arg.Subj,
		arg.CreatedAt,
	)
	return err
}
