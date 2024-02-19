// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: follows.sql

package store_queries

import (
	"context"
	"database/sql"
	"time"
)

const countFollowersByTarget = `-- name: CountFollowersByTarget :one
SELECT COUNT(*)
FROM follows
WHERE target_did = $1
`

func (q *Queries) CountFollowersByTarget(ctx context.Context, targetDid string) (int64, error) {
	row := q.queryRow(ctx, q.countFollowersByTargetStmt, countFollowersByTarget, targetDid)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countFollowsByActor = `-- name: CountFollowsByActor :one
SELECT COUNT(*)
FROM follows
WHERE actor_did = $1
`

func (q *Queries) CountFollowsByActor(ctx context.Context, actorDid string) (int64, error) {
	row := q.queryRow(ctx, q.countFollowsByActorStmt, countFollowsByActor, actorDid)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countFollowsByActorAndTarget = `-- name: CountFollowsByActorAndTarget :one
SELECT COUNT(*)
FROM follows
WHERE actor_did = $1
    AND target_did = $2
`

type CountFollowsByActorAndTargetParams struct {
	ActorDid  string `json:"actor_did"`
	TargetDid string `json:"target_did"`
}

func (q *Queries) CountFollowsByActorAndTarget(ctx context.Context, arg CountFollowsByActorAndTargetParams) (int64, error) {
	row := q.queryRow(ctx, q.countFollowsByActorAndTargetStmt, countFollowsByActorAndTarget, arg.ActorDid, arg.TargetDid)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createFollow = `-- name: CreateFollow :exec
INSERT INTO follows(
        actor_did,
        rkey,
        target_did,
        created_at
    )
VALUES ($1, $2, $3, $4)
`

type CreateFollowParams struct {
	ActorDid  string       `json:"actor_did"`
	Rkey      string       `json:"rkey"`
	TargetDid string       `json:"target_did"`
	CreatedAt sql.NullTime `json:"created_at"`
}

func (q *Queries) CreateFollow(ctx context.Context, arg CreateFollowParams) error {
	_, err := q.exec(ctx, q.createFollowStmt, createFollow,
		arg.ActorDid,
		arg.Rkey,
		arg.TargetDid,
		arg.CreatedAt,
	)
	return err
}

const deleteFollow = `-- name: DeleteFollow :exec
DELETE FROM follows
WHERE actor_did = $1
    AND rkey = $2
`

type DeleteFollowParams struct {
	ActorDid string `json:"actor_did"`
	Rkey     string `json:"rkey"`
}

func (q *Queries) DeleteFollow(ctx context.Context, arg DeleteFollowParams) error {
	_, err := q.exec(ctx, q.deleteFollowStmt, deleteFollow, arg.ActorDid, arg.Rkey)
	return err
}

const deleteFollowsByActor = `-- name: DeleteFollowsByActor :exec
DELETE FROM follows
WHERE actor_did = $1
`

func (q *Queries) DeleteFollowsByActor(ctx context.Context, actorDid string) error {
	_, err := q.exec(ctx, q.deleteFollowsByActorStmt, deleteFollowsByActor, actorDid)
	return err
}

const deleteFollowsByTarget = `-- name: DeleteFollowsByTarget :exec
DELETE FROM follows
WHERE target_did = $1
`

func (q *Queries) DeleteFollowsByTarget(ctx context.Context, targetDid string) error {
	_, err := q.exec(ctx, q.deleteFollowsByTargetStmt, deleteFollowsByTarget, targetDid)
	return err
}

const getFollow = `-- name: GetFollow :one
SELECT actor_did, rkey, target_did, created_at, inserted_at
FROM follows
WHERE actor_did = $1
    AND rkey = $2
`

type GetFollowParams struct {
	ActorDid string `json:"actor_did"`
	Rkey     string `json:"rkey"`
}

func (q *Queries) GetFollow(ctx context.Context, arg GetFollowParams) (Follow, error) {
	row := q.queryRow(ctx, q.getFollowStmt, getFollow, arg.ActorDid, arg.Rkey)
	var i Follow
	err := row.Scan(
		&i.ActorDid,
		&i.Rkey,
		&i.TargetDid,
		&i.CreatedAt,
		&i.InsertedAt,
	)
	return i, err
}

const getFollowPage = `-- name: GetFollowPage :many
SELECT actor_did, rkey, target_did, created_at, inserted_at
FROM follows
WHERE inserted_at > $1
ORDER BY inserted_at
LIMIT $2
`

type GetFollowPageParams struct {
	InsertedAt time.Time `json:"inserted_at"`
	Limit      int32     `json:"limit"`
}

func (q *Queries) GetFollowPage(ctx context.Context, arg GetFollowPageParams) ([]Follow, error) {
	rows, err := q.query(ctx, q.getFollowPageStmt, getFollowPage, arg.InsertedAt, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Follow
	for rows.Next() {
		var i Follow
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.TargetDid,
			&i.CreatedAt,
			&i.InsertedAt,
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

const getFollowsByActor = `-- name: GetFollowsByActor :many
SELECT actor_did, rkey, target_did, created_at, inserted_at
FROM follows
WHERE actor_did = $1
ORDER BY created_at DESC
LIMIT $2
`

type GetFollowsByActorParams struct {
	ActorDid string `json:"actor_did"`
	Limit    int32  `json:"limit"`
}

func (q *Queries) GetFollowsByActor(ctx context.Context, arg GetFollowsByActorParams) ([]Follow, error) {
	rows, err := q.query(ctx, q.getFollowsByActorStmt, getFollowsByActor, arg.ActorDid, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Follow
	for rows.Next() {
		var i Follow
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.TargetDid,
			&i.CreatedAt,
			&i.InsertedAt,
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

const getFollowsByActorAndTarget = `-- name: GetFollowsByActorAndTarget :many
SELECT actor_did, rkey, target_did, created_at, inserted_at
FROM follows
WHERE actor_did = $1
    AND target_did = $2
ORDER BY created_at DESC
LIMIT $3
`

type GetFollowsByActorAndTargetParams struct {
	ActorDid  string `json:"actor_did"`
	TargetDid string `json:"target_did"`
	Limit     int32  `json:"limit"`
}

func (q *Queries) GetFollowsByActorAndTarget(ctx context.Context, arg GetFollowsByActorAndTargetParams) ([]Follow, error) {
	rows, err := q.query(ctx, q.getFollowsByActorAndTargetStmt, getFollowsByActorAndTarget, arg.ActorDid, arg.TargetDid, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Follow
	for rows.Next() {
		var i Follow
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.TargetDid,
			&i.CreatedAt,
			&i.InsertedAt,
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

const getFollowsByTarget = `-- name: GetFollowsByTarget :many
SELECT actor_did, rkey, target_did, created_at, inserted_at
FROM follows
WHERE target_did = $1
ORDER BY created_at DESC
LIMIT $2
`

type GetFollowsByTargetParams struct {
	TargetDid string `json:"target_did"`
	Limit     int32  `json:"limit"`
}

func (q *Queries) GetFollowsByTarget(ctx context.Context, arg GetFollowsByTargetParams) ([]Follow, error) {
	rows, err := q.query(ctx, q.getFollowsByTargetStmt, getFollowsByTarget, arg.TargetDid, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Follow
	for rows.Next() {
		var i Follow
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.TargetDid,
			&i.CreatedAt,
			&i.InsertedAt,
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
