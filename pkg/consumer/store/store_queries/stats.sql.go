// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: stats.sql

package store_queries

import (
	"context"
)

const getDailySummaries = `-- name: GetDailySummaries :many
SELECT date, "Likes per Day", "Daily Active Likers", "Daily Active Posters", "Posts per Day", "Posts with Images per Day", "Images per Day", "Images with Alt Text per Day", "First Time Posters", "Follows per Day", "Daily Active Followers", "Blocks per Day", "Daily Active Blockers"
FROM daily_summary
`

func (q *Queries) GetDailySummaries(ctx context.Context) ([]DailySummary, error) {
	rows, err := q.query(ctx, q.getDailySummariesStmt, getDailySummaries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DailySummary
	for rows.Next() {
		var i DailySummary
		if err := rows.Scan(
			&i.Date,
			&i.LikesPerDay,
			&i.DailyActiveLikers,
			&i.DailyActivePosters,
			&i.PostsPerDay,
			&i.PostsWithImagesPerDay,
			&i.ImagesPerDay,
			&i.ImagesWithAltTextPerDay,
			&i.FirstTimePosters,
			&i.FollowsPerDay,
			&i.DailyActiveFollowers,
			&i.BlocksPerDay,
			&i.DailyActiveBlockers,
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

const getFollowerPercentiles = `-- name: GetFollowerPercentiles :one
WITH all_counts AS (
    SELECT a.did,
        COALESCE(fc.num_followers, 0) AS followers
    FROM actors a
        LEFT JOIN follower_counts fc ON a.did = fc.actor_did
),
percentiles AS (
    SELECT percentile_cont(
            ARRAY [0.25, 0.50, 0.75, 0.90, 0.95, 0.99, 0.995, 0.997, 0.999, 0.9999]
        ) WITHIN GROUP (
            ORDER BY followers
        ) AS pct
    FROM all_counts
)
SELECT pct [1]::float AS p25,
    pct [2]::float AS p50,
    pct [3]::float AS p75,
    pct [4]::float AS p90,
    pct [5]::float AS p95,
    pct [6]::float AS p99,
    pct [7]::float AS p99_5,
    pct [8]::float AS p99_7,
    pct [9]::float AS p99_9,
    pct [10]::float AS p99_99
FROM percentiles
`

type GetFollowerPercentilesRow struct {
	P25   float64 `json:"p25"`
	P50   float64 `json:"p50"`
	P75   float64 `json:"p75"`
	P90   float64 `json:"p90"`
	P95   float64 `json:"p95"`
	P99   float64 `json:"p99"`
	P995  float64 `json:"p99_5"`
	P997  float64 `json:"p99_7"`
	P999  float64 `json:"p99_9"`
	P9999 float64 `json:"p99_99"`
}

func (q *Queries) GetFollowerPercentiles(ctx context.Context) (GetFollowerPercentilesRow, error) {
	row := q.queryRow(ctx, q.getFollowerPercentilesStmt, getFollowerPercentiles)
	var i GetFollowerPercentilesRow
	err := row.Scan(
		&i.P25,
		&i.P50,
		&i.P75,
		&i.P90,
		&i.P95,
		&i.P99,
		&i.P995,
		&i.P997,
		&i.P999,
		&i.P9999,
	)
	return i, err
}
