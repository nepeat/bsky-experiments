// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: stats.sql

package store_queries

import (
	"context"
)

const getDailySummaries = `-- name: GetDailySummaries :many
SELECT date, "Likes per Day", "Daily Active Likers", "Daily Active Posters", "Posts per Day", "Posts with Images per Day", "Images per Day", "Images with Alt Text per Day", "First Time Posters", "Follows per Day", "Daily Active Followers", "Blocks per Day", "Daily Active Blockers"
FROM daily_summary
ORDER BY date DESC
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
SELECT p25, p50, p75, p90, p95, p99, p99_5, p99_7, p99_9, p99_99
FROM follower_stats
`

func (q *Queries) GetFollowerPercentiles(ctx context.Context) (FollowerStat, error) {
	row := q.queryRow(ctx, q.getFollowerPercentilesStmt, getFollowerPercentiles)
	var i FollowerStat
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
