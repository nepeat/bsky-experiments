// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: recent_posts.sql

package store_queries

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/sqlc-dev/pqtype"
)

const createRecentPost = `-- name: CreateRecentPost :exec
INSERT INTO recent_posts (
        actor_did,
        rkey,
        content,
        parent_post_actor_did,
        parent_post_rkey,
        quote_post_actor_did,
        quote_post_rkey,
        root_post_actor_did,
        root_post_rkey,
        has_embedded_media,
        facets,
        embed,
        tags,
        created_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13,
        $14
    )
`

type CreateRecentPostParams struct {
	ActorDid           string                `json:"actor_did"`
	Rkey               string                `json:"rkey"`
	Content            sql.NullString        `json:"content"`
	ParentPostActorDid sql.NullString        `json:"parent_post_actor_did"`
	ParentPostRkey     sql.NullString        `json:"parent_post_rkey"`
	QuotePostActorDid  sql.NullString        `json:"quote_post_actor_did"`
	QuotePostRkey      sql.NullString        `json:"quote_post_rkey"`
	RootPostActorDid   sql.NullString        `json:"root_post_actor_did"`
	RootPostRkey       sql.NullString        `json:"root_post_rkey"`
	HasEmbeddedMedia   bool                  `json:"has_embedded_media"`
	Facets             pqtype.NullRawMessage `json:"facets"`
	Embed              pqtype.NullRawMessage `json:"embed"`
	Tags               []string              `json:"tags"`
	CreatedAt          sql.NullTime          `json:"created_at"`
}

func (q *Queries) CreateRecentPost(ctx context.Context, arg CreateRecentPostParams) error {
	_, err := q.exec(ctx, q.createRecentPostStmt, createRecentPost,
		arg.ActorDid,
		arg.Rkey,
		arg.Content,
		arg.ParentPostActorDid,
		arg.ParentPostRkey,
		arg.QuotePostActorDid,
		arg.QuotePostRkey,
		arg.RootPostActorDid,
		arg.RootPostRkey,
		arg.HasEmbeddedMedia,
		arg.Facets,
		arg.Embed,
		pq.Array(arg.Tags),
		arg.CreatedAt,
	)
	return err
}

const deleteRecentPost = `-- name: DeleteRecentPost :exec
DELETE FROM recent_posts
WHERE actor_did = $1
    AND rkey = $2
`

type DeleteRecentPostParams struct {
	ActorDid string `json:"actor_did"`
	Rkey     string `json:"rkey"`
}

func (q *Queries) DeleteRecentPost(ctx context.Context, arg DeleteRecentPostParams) error {
	_, err := q.exec(ctx, q.deleteRecentPostStmt, deleteRecentPost, arg.ActorDid, arg.Rkey)
	return err
}

const getRecentPost = `-- name: GetRecentPost :one
SELECT actor_did, rkey, content, parent_post_actor_did, quote_post_actor_did, quote_post_rkey, parent_post_rkey, root_post_actor_did, root_post_rkey, facets, embed, tags, has_embedded_media, created_at, inserted_at
FROM recent_posts
WHERE actor_did = $1
    AND rkey = $2
`

type GetRecentPostParams struct {
	ActorDid string `json:"actor_did"`
	Rkey     string `json:"rkey"`
}

func (q *Queries) GetRecentPost(ctx context.Context, arg GetRecentPostParams) (RecentPost, error) {
	row := q.queryRow(ctx, q.getRecentPostStmt, getRecentPost, arg.ActorDid, arg.Rkey)
	var i RecentPost
	err := row.Scan(
		&i.ActorDid,
		&i.Rkey,
		&i.Content,
		&i.ParentPostActorDid,
		&i.QuotePostActorDid,
		&i.QuotePostRkey,
		&i.ParentPostRkey,
		&i.RootPostActorDid,
		&i.RootPostRkey,
		&i.Facets,
		&i.Embed,
		pq.Array(&i.Tags),
		&i.HasEmbeddedMedia,
		&i.CreatedAt,
		&i.InsertedAt,
	)
	return i, err
}

const getRecentPostsByActor = `-- name: GetRecentPostsByActor :many
SELECT actor_did, rkey, content, parent_post_actor_did, quote_post_actor_did, quote_post_rkey, parent_post_rkey, root_post_actor_did, root_post_rkey, facets, embed, tags, has_embedded_media, created_at, inserted_at
FROM recent_posts
WHERE actor_did = $1
ORDER BY created_at DESC
LIMIT $2
`

type GetRecentPostsByActorParams struct {
	ActorDid string `json:"actor_did"`
	Limit    int32  `json:"limit"`
}

func (q *Queries) GetRecentPostsByActor(ctx context.Context, arg GetRecentPostsByActorParams) ([]RecentPost, error) {
	rows, err := q.query(ctx, q.getRecentPostsByActorStmt, getRecentPostsByActor, arg.ActorDid, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RecentPost
	for rows.Next() {
		var i RecentPost
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.Content,
			&i.ParentPostActorDid,
			&i.QuotePostActorDid,
			&i.QuotePostRkey,
			&i.ParentPostRkey,
			&i.RootPostActorDid,
			&i.RootPostRkey,
			&i.Facets,
			&i.Embed,
			pq.Array(&i.Tags),
			&i.HasEmbeddedMedia,
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

const getRecentPostsByActorsFollowingTarget = `-- name: GetRecentPostsByActorsFollowingTarget :many
WITH followers AS (
    SELECT actor_did
    FROM follows
    WHERE target_did = $1
)
SELECT p.actor_did, p.rkey, p.content, p.parent_post_actor_did, p.quote_post_actor_did, p.quote_post_rkey, p.parent_post_rkey, p.root_post_actor_did, p.root_post_rkey, p.facets, p.embed, p.tags, p.has_embedded_media, p.created_at, p.inserted_at
FROM recent_posts p
    JOIN followers f ON f.actor_did = p.actor_did
WHERE (p.created_at, p.actor_did, p.rkey) < (
        $3::TIMESTAMPTZ,
        $4::TEXT,
        $5::TEXT
    )
    AND (p.root_post_rkey IS NULL)
    AND (
        (p.parent_relationship IS NULL)
        OR (p.parent_relationship <> 'r'::text)
    )
ORDER BY p.created_at DESC,
    p.actor_did DESC,
    p.rkey DESC
LIMIT $2
`

type GetRecentPostsByActorsFollowingTargetParams struct {
	TargetDid       string    `json:"target_did"`
	Limit           int32     `json:"limit"`
	CursorCreatedAt time.Time `json:"cursor_created_at"`
	CursorActorDid  string    `json:"cursor_actor_did"`
	CursorRkey      string    `json:"cursor_rkey"`
}

func (q *Queries) GetRecentPostsByActorsFollowingTarget(ctx context.Context, arg GetRecentPostsByActorsFollowingTargetParams) ([]RecentPost, error) {
	rows, err := q.query(ctx, q.getRecentPostsByActorsFollowingTargetStmt, getRecentPostsByActorsFollowingTarget,
		arg.TargetDid,
		arg.Limit,
		arg.CursorCreatedAt,
		arg.CursorActorDid,
		arg.CursorRkey,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RecentPost
	for rows.Next() {
		var i RecentPost
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.Content,
			&i.ParentPostActorDid,
			&i.QuotePostActorDid,
			&i.QuotePostRkey,
			&i.ParentPostRkey,
			&i.RootPostActorDid,
			&i.RootPostRkey,
			&i.Facets,
			&i.Embed,
			pq.Array(&i.Tags),
			&i.HasEmbeddedMedia,
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

const getRecentPostsFromNonMoots = `-- name: GetRecentPostsFromNonMoots :many
WITH my_follows AS (
    SELECT target_did
    FROM follows
    WHERE follows.actor_did = $1
),
non_moots AS (
    SELECT actor_did
    FROM follows f
        LEFT JOIN my_follows ON f.actor_did = my_follows.target_did
    WHERE f.target_did = $1
        AND my_follows.target_did IS NULL
),
non_moots_and_non_spam AS (
    SELECT nm.actor_did
    FROM non_moots nm
        LEFT JOIN following_counts fc ON nm.actor_did = fc.actor_did
    WHERE fc.num_following < 4000
)
SELECT p.actor_did, p.rkey, p.content, p.parent_post_actor_did, p.quote_post_actor_did, p.quote_post_rkey, p.parent_post_rkey, p.root_post_actor_did, p.root_post_rkey, p.facets, p.embed, p.tags, p.has_embedded_media, p.created_at, p.inserted_at
FROM recent_posts p
    JOIN non_moots_and_non_spam f ON f.actor_did = p.actor_did
WHERE (p.created_at, p.actor_did, p.rkey) < (
        $3::TIMESTAMPTZ,
        $4::TEXT,
        $5::TEXT
    )
    AND p.root_post_rkey IS NULL
    AND p.parent_post_rkey IS NULL
    AND p.created_at > NOW() - make_interval(hours := 24)
ORDER BY p.created_at DESC,
    p.actor_did DESC,
    p.rkey DESC
LIMIT $2
`

type GetRecentPostsFromNonMootsParams struct {
	ActorDid        string    `json:"actor_did"`
	Limit           int32     `json:"limit"`
	CursorCreatedAt time.Time `json:"cursor_created_at"`
	CursorActorDid  string    `json:"cursor_actor_did"`
	CursorRkey      string    `json:"cursor_rkey"`
}

func (q *Queries) GetRecentPostsFromNonMoots(ctx context.Context, arg GetRecentPostsFromNonMootsParams) ([]RecentPost, error) {
	rows, err := q.query(ctx, q.getRecentPostsFromNonMootsStmt, getRecentPostsFromNonMoots,
		arg.ActorDid,
		arg.Limit,
		arg.CursorCreatedAt,
		arg.CursorActorDid,
		arg.CursorRkey,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RecentPost
	for rows.Next() {
		var i RecentPost
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.Content,
			&i.ParentPostActorDid,
			&i.QuotePostActorDid,
			&i.QuotePostRkey,
			&i.ParentPostRkey,
			&i.RootPostActorDid,
			&i.RootPostRkey,
			&i.Facets,
			&i.Embed,
			pq.Array(&i.Tags),
			&i.HasEmbeddedMedia,
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

const getRecentPostsFromNonSpamUsers = `-- name: GetRecentPostsFromNonSpamUsers :many
WITH non_spam AS (
    SELECT nm.actor_did
    FROM unnest($5::TEXT []) nm(actor_did)
        LEFT JOIN following_counts fc ON nm.actor_did = fc.actor_did
    WHERE fc.num_following < 4000
)
SELECT p.actor_did, p.rkey, p.content, p.parent_post_actor_did, p.quote_post_actor_did, p.quote_post_rkey, p.parent_post_rkey, p.root_post_actor_did, p.root_post_rkey, p.facets, p.embed, p.tags, p.has_embedded_media, p.created_at, p.inserted_at
FROM recent_posts p
    JOIN non_spam f ON f.actor_did = p.actor_did
WHERE (p.created_at, p.actor_did, p.rkey) < (
        $2::TIMESTAMPTZ,
        $3::TEXT,
        $4::TEXT
    )
    AND p.root_post_rkey IS NULL
    AND p.parent_post_rkey IS NULL
    AND p.created_at > NOW() - make_interval(hours := 24)
ORDER BY p.created_at DESC,
    p.actor_did DESC,
    p.rkey DESC
LIMIT $1
`

type GetRecentPostsFromNonSpamUsersParams struct {
	Limit           int32     `json:"limit"`
	CursorCreatedAt time.Time `json:"cursor_created_at"`
	CursorActorDid  string    `json:"cursor_actor_did"`
	CursorRkey      string    `json:"cursor_rkey"`
	Dids            []string  `json:"dids"`
}

func (q *Queries) GetRecentPostsFromNonSpamUsers(ctx context.Context, arg GetRecentPostsFromNonSpamUsersParams) ([]RecentPost, error) {
	rows, err := q.query(ctx, q.getRecentPostsFromNonSpamUsersStmt, getRecentPostsFromNonSpamUsers,
		arg.Limit,
		arg.CursorCreatedAt,
		arg.CursorActorDid,
		arg.CursorRkey,
		pq.Array(arg.Dids),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RecentPost
	for rows.Next() {
		var i RecentPost
		if err := rows.Scan(
			&i.ActorDid,
			&i.Rkey,
			&i.Content,
			&i.ParentPostActorDid,
			&i.QuotePostActorDid,
			&i.QuotePostRkey,
			&i.ParentPostRkey,
			&i.RootPostActorDid,
			&i.RootPostRkey,
			&i.Facets,
			&i.Embed,
			pq.Array(&i.Tags),
			&i.HasEmbeddedMedia,
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

const trimOldRecentPosts = `-- name: TrimOldRecentPosts :execrows
DELETE FROM recent_posts
WHERE created_at < NOW() - make_interval(hours := $1)
    AND created_at > NOW() + make_interval(mins := 15)
`

func (q *Queries) TrimOldRecentPosts(ctx context.Context, hours int32) (int64, error) {
	result, err := q.exec(ctx, q.trimOldRecentPostsStmt, trimOldRecentPosts, hours)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
