// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: points.sql

package store_queries

import (
	"context"
)

const createPointAssignment = `-- name: CreatePointAssignment :exec
INSERT INTO point_assignments (
        event_id,
        actor_did,
        points
    )
VALUES ($1, $2, $3) ON CONFLICT (event_id, actor_did) DO
UPDATE
SET points = EXCLUDED.points
`

type CreatePointAssignmentParams struct {
	EventID  int64  `json:"event_id"`
	ActorDid string `json:"actor_did"`
	Points   int32  `json:"points"`
}

// CREATE TABLE point_assignments (
//
//	id BIGSERIAL PRIMARY KEY,
//	event_id BIGINT NOT NULL,
//	actor_did TEXT NOT NULL,
//	points INTEGER NOT NULL,
//	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
//	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
//
// );
func (q *Queries) CreatePointAssignment(ctx context.Context, arg CreatePointAssignmentParams) error {
	_, err := q.exec(ctx, q.createPointAssignmentStmt, createPointAssignment, arg.EventID, arg.ActorDid, arg.Points)
	return err
}

const deletePointAssignment = `-- name: DeletePointAssignment :exec
DELETE FROM point_assignments
WHERE event_id = $1
    AND actor_did = $2
`

type DeletePointAssignmentParams struct {
	EventID  int64  `json:"event_id"`
	ActorDid string `json:"actor_did"`
}

func (q *Queries) DeletePointAssignment(ctx context.Context, arg DeletePointAssignmentParams) error {
	_, err := q.exec(ctx, q.deletePointAssignmentStmt, deletePointAssignment, arg.EventID, arg.ActorDid)
	return err
}

const getPointAssignment = `-- name: GetPointAssignment :one
SELECT id, event_id, actor_did, points, created_at, updated_at
FROM point_assignments
WHERE event_id = $1
    AND actor_did = $2
LIMIT 1
`

type GetPointAssignmentParams struct {
	EventID  int64  `json:"event_id"`
	ActorDid string `json:"actor_did"`
}

func (q *Queries) GetPointAssignment(ctx context.Context, arg GetPointAssignmentParams) (PointAssignment, error) {
	row := q.queryRow(ctx, q.getPointAssignmentStmt, getPointAssignment, arg.EventID, arg.ActorDid)
	var i PointAssignment
	err := row.Scan(
		&i.ID,
		&i.EventID,
		&i.ActorDid,
		&i.Points,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPointAssignmentsForActor = `-- name: GetPointAssignmentsForActor :many
SELECT id, event_id, actor_did, points, created_at, updated_at
FROM point_assignments
WHERE actor_did = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type GetPointAssignmentsForActorParams struct {
	ActorDid string `json:"actor_did"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

func (q *Queries) GetPointAssignmentsForActor(ctx context.Context, arg GetPointAssignmentsForActorParams) ([]PointAssignment, error) {
	rows, err := q.query(ctx, q.getPointAssignmentsForActorStmt, getPointAssignmentsForActor, arg.ActorDid, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PointAssignment
	for rows.Next() {
		var i PointAssignment
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.ActorDid,
			&i.Points,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getPointAssignmentsForEvent = `-- name: GetPointAssignmentsForEvent :many
SELECT id, event_id, actor_did, points, created_at, updated_at
FROM point_assignments
WHERE event_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type GetPointAssignmentsForEventParams struct {
	EventID int64 `json:"event_id"`
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
}

func (q *Queries) GetPointAssignmentsForEvent(ctx context.Context, arg GetPointAssignmentsForEventParams) ([]PointAssignment, error) {
	rows, err := q.query(ctx, q.getPointAssignmentsForEventStmt, getPointAssignmentsForEvent, arg.EventID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PointAssignment
	for rows.Next() {
		var i PointAssignment
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.ActorDid,
			&i.Points,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getTotalPointsForActor = `-- name: GetTotalPointsForActor :one
SELECT SUM(points)
FROM point_assignments
WHERE actor_did = $1
`

func (q *Queries) GetTotalPointsForActor(ctx context.Context, actorDid string) (int64, error) {
	row := q.queryRow(ctx, q.getTotalPointsForActorStmt, getTotalPointsForActor, actorDid)
	var sum int64
	err := row.Scan(&sum)
	return sum, err
}

const getTotalPointsForEvent = `-- name: GetTotalPointsForEvent :one
SELECT SUM(points)
FROM point_assignments
WHERE event_id = $1
`

func (q *Queries) GetTotalPointsForEvent(ctx context.Context, eventID int64) (int64, error) {
	row := q.queryRow(ctx, q.getTotalPointsForEventStmt, getTotalPointsForEvent, eventID)
	var sum int64
	err := row.Scan(&sum)
	return sum, err
}

const upatePointAssignment = `-- name: UpatePointAssignment :exec
UPDATE point_assignments
SET points = $3
WHERE event_id = $1
    AND actor_did = $2
`

type UpatePointAssignmentParams struct {
	EventID  int64  `json:"event_id"`
	ActorDid string `json:"actor_did"`
	Points   int32  `json:"points"`
}

func (q *Queries) UpatePointAssignment(ctx context.Context, arg UpatePointAssignmentParams) error {
	_, err := q.exec(ctx, q.upatePointAssignmentStmt, upatePointAssignment, arg.EventID, arg.ActorDid, arg.Points)
	return err
}