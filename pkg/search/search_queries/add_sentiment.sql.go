// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: add_sentiment.sql

package search_queries

import (
	"context"
	"database/sql"
)

const setPostSentiment = `-- name: SetPostSentiment :exec
UPDATE posts
SET sentiment = $1,
    sentiment_confidence = $2
WHERE id = $3
    AND author_did = $4
`

type SetPostSentimentParams struct {
	Sentiment           sql.NullString  `json:"sentiment"`
	SentimentConfidence sql.NullFloat64 `json:"sentiment_confidence"`
	ID                  string          `json:"id"`
	AuthorDid           string          `json:"author_did"`
}

func (q *Queries) SetPostSentiment(ctx context.Context, arg SetPostSentimentParams) error {
	_, err := q.exec(ctx, q.setPostSentimentStmt, setPostSentiment,
		arg.Sentiment,
		arg.SentimentConfidence,
		arg.ID,
		arg.AuthorDid,
	)
	return err
}
