// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: follow_counts.sql

package store_queries

import (
	"context"
	"time"
)

const decrementFollowerCountByN = `-- name: DecrementFollowerCountByN :exec
INSERT INTO follower_counts (actor_did, num_followers, updated_at)
SELECT $1,
    - $3::int,
    $2 ON CONFLICT (actor_did) DO
UPDATE
SET (num_followers, updated_at) = (
        GREATEST(
            0,
            follower_counts.num_followers + EXCLUDED.num_followers
        ),
        EXCLUDED.updated_at
    )
`

type DecrementFollowerCountByNParams struct {
	ActorDid     string    `json:"actor_did"`
	UpdatedAt    time.Time `json:"updated_at"`
	NumFollowers int32     `json:"num_followers"`
}

func (q *Queries) DecrementFollowerCountByN(ctx context.Context, arg DecrementFollowerCountByNParams) error {
	_, err := q.exec(ctx, q.decrementFollowerCountByNStmt, decrementFollowerCountByN, arg.ActorDid, arg.UpdatedAt, arg.NumFollowers)
	return err
}

const decrementFollowingCountByN = `-- name: DecrementFollowingCountByN :exec
INSERT INTO following_counts (actor_did, num_following, updated_at)
SELECT $1,
    - $3::int,
    $2 ON CONFLICT (actor_did) DO
UPDATE
SET (num_following, updated_at) = (
        GREATEST(
            0,
            following_counts.num_following + EXCLUDED.num_following
        ),
        EXCLUDED.updated_at
    )
`

type DecrementFollowingCountByNParams struct {
	ActorDid     string    `json:"actor_did"`
	UpdatedAt    time.Time `json:"updated_at"`
	NumFollowing int32     `json:"num_following"`
}

func (q *Queries) DecrementFollowingCountByN(ctx context.Context, arg DecrementFollowingCountByNParams) error {
	_, err := q.exec(ctx, q.decrementFollowingCountByNStmt, decrementFollowingCountByN, arg.ActorDid, arg.UpdatedAt, arg.NumFollowing)
	return err
}

const deleteFollowerCount = `-- name: DeleteFollowerCount :exec
DELETE FROM follower_counts
WHERE actor_did = $1
`

func (q *Queries) DeleteFollowerCount(ctx context.Context, actorDid string) error {
	_, err := q.exec(ctx, q.deleteFollowerCountStmt, deleteFollowerCount, actorDid)
	return err
}

const deleteFollowingCount = `-- name: DeleteFollowingCount :exec
DELETE FROM following_counts
WHERE actor_did = $1
`

func (q *Queries) DeleteFollowingCount(ctx context.Context, actorDid string) error {
	_, err := q.exec(ctx, q.deleteFollowingCountStmt, deleteFollowingCount, actorDid)
	return err
}

const getFollowerCount = `-- name: GetFollowerCount :one
SELECT nfc.actor_did, nfc.num_followers, nfc.updated_at
FROM follower_counts nfc
WHERE nfc.actor_did = $1
`

func (q *Queries) GetFollowerCount(ctx context.Context, actorDid string) (FollowerCount, error) {
	row := q.queryRow(ctx, q.getFollowerCountStmt, getFollowerCount, actorDid)
	var i FollowerCount
	err := row.Scan(&i.ActorDid, &i.NumFollowers, &i.UpdatedAt)
	return i, err
}

const getFollowingCount = `-- name: GetFollowingCount :one
SELECT nfc.actor_did, nfc.num_following, nfc.updated_at
FROM following_counts nfc
WHERE nfc.actor_did = $1
`

func (q *Queries) GetFollowingCount(ctx context.Context, actorDid string) (FollowingCount, error) {
	row := q.queryRow(ctx, q.getFollowingCountStmt, getFollowingCount, actorDid)
	var i FollowingCount
	err := row.Scan(&i.ActorDid, &i.NumFollowing, &i.UpdatedAt)
	return i, err
}

const incrementFollowerCountByN = `-- name: IncrementFollowerCountByN :exec
INSERT INTO follower_counts (actor_did, num_followers, updated_at)
SELECT $1,
    $3::int,
    $2 ON CONFLICT (actor_did) DO
UPDATE
SET (num_followers, updated_at) = (
        follower_counts.num_followers + EXCLUDED.num_followers,
        EXCLUDED.updated_at
    )
`

type IncrementFollowerCountByNParams struct {
	ActorDid     string    `json:"actor_did"`
	UpdatedAt    time.Time `json:"updated_at"`
	NumFollowers int32     `json:"num_followers"`
}

// Follow Counts
func (q *Queries) IncrementFollowerCountByN(ctx context.Context, arg IncrementFollowerCountByNParams) error {
	_, err := q.exec(ctx, q.incrementFollowerCountByNStmt, incrementFollowerCountByN, arg.ActorDid, arg.UpdatedAt, arg.NumFollowers)
	return err
}

const incrementFollowingCountByN = `-- name: IncrementFollowingCountByN :exec
INSERT INTO following_counts (actor_did, num_following, updated_at)
SELECT $1,
    $3::int,
    $2 ON CONFLICT (actor_did) DO
UPDATE
SET (num_following, updated_at) = (
        following_counts.num_following + EXCLUDED.num_following,
        EXCLUDED.updated_at
    )
`

type IncrementFollowingCountByNParams struct {
	ActorDid     string    `json:"actor_did"`
	UpdatedAt    time.Time `json:"updated_at"`
	NumFollowing int32     `json:"num_following"`
}

// Following Counts
func (q *Queries) IncrementFollowingCountByN(ctx context.Context, arg IncrementFollowingCountByNParams) error {
	_, err := q.exec(ctx, q.incrementFollowingCountByNStmt, incrementFollowingCountByN, arg.ActorDid, arg.UpdatedAt, arg.NumFollowing)
	return err
}
