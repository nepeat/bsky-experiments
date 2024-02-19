// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package search_queries

import (
	"database/sql"
	"time"

	"github.com/sqlc-dev/pqtype"
)

type Author struct {
	Did           string `json:"did"`
	Handle        string `json:"handle"`
	ClusterOptOut bool   `json:"cluster_opt_out"`
}

type AuthorBlock struct {
	ActorDid  string    `json:"actor_did"`
	TargetDid string    `json:"target_did"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthorCluster struct {
	AuthorDid string `json:"author_did"`
	ClusterID int32  `json:"cluster_id"`
}

type AuthorLabel struct {
	AuthorDid string    `json:"author_did"`
	LabelID   int32     `json:"label_id"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthorStat struct {
	TotalAuthors       int64   `json:"total_authors"`
	TotalPosts         int64   `json:"total_posts"`
	MeanPostsPerAuthor float64 `json:"mean_posts_per_author"`
	Gt1                int64   `json:"gt_1"`
	Gt5                int64   `json:"gt_5"`
	Gt10               int64   `json:"gt_10"`
	Gt20               int64   `json:"gt_20"`
	Gt100              int64   `json:"gt_100"`
	Gt1000             int64   `json:"gt_1000"`
	P25                int64   `json:"p25"`
	P50                int64   `json:"p50"`
	P75                int64   `json:"p75"`
	P90                int64   `json:"p90"`
	P95                int64   `json:"p95"`
	P99                int64   `json:"p99"`
	P995               int64   `json:"p99_5"`
	P997               int64   `json:"p99_7"`
	P999               int64   `json:"p99_9"`
	P9999              int64   `json:"p99_99"`
}

type Cluster struct {
	ID          int32  `json:"id"`
	LookupAlias string `json:"lookup_alias"`
	Name        string `json:"name"`
}

type Image struct {
	Cid         string                `json:"cid"`
	PostID      string                `json:"post_id"`
	AuthorDid   string                `json:"author_did"`
	AltText     sql.NullString        `json:"alt_text"`
	MimeType    string                `json:"mime_type"`
	CreatedAt   time.Time             `json:"created_at"`
	CvCompleted bool                  `json:"cv_completed"`
	CvRunAt     sql.NullTime          `json:"cv_run_at"`
	CvClasses   pqtype.NullRawMessage `json:"cv_classes"`
}

type Label struct {
	ID          int64  `json:"id"`
	LookupAlias string `json:"lookup_alias"`
	Name        string `json:"name"`
}

type Post struct {
	ID                  string          `json:"id"`
	Text                string          `json:"text"`
	ParentPostID        sql.NullString  `json:"parent_post_id"`
	RootPostID          sql.NullString  `json:"root_post_id"`
	AuthorDid           string          `json:"author_did"`
	CreatedAt           time.Time       `json:"created_at"`
	HasEmbeddedMedia    bool            `json:"has_embedded_media"`
	ParentRelationship  sql.NullString  `json:"parent_relationship"`
	Sentiment           sql.NullString  `json:"sentiment"`
	SentimentConfidence sql.NullFloat64 `json:"sentiment_confidence"`
	IndexedAt           sql.NullTime    `json:"indexed_at"`
}

type PostHotness struct {
	ID                  string          `json:"id"`
	Text                string          `json:"text"`
	ParentPostID        sql.NullString  `json:"parent_post_id"`
	RootPostID          sql.NullString  `json:"root_post_id"`
	AuthorDid           string          `json:"author_did"`
	CreatedAt           time.Time       `json:"created_at"`
	HasEmbeddedMedia    bool            `json:"has_embedded_media"`
	ParentRelationship  sql.NullString  `json:"parent_relationship"`
	Sentiment           sql.NullString  `json:"sentiment"`
	SentimentConfidence sql.NullFloat64 `json:"sentiment_confidence"`
	PostLabels          interface{}     `json:"post_labels"`
	ClusterLabel        sql.NullString  `json:"cluster_label"`
	AuthorLabels        interface{}     `json:"author_labels"`
	Hotness             float64         `json:"hotness"`
}

type PostLabel struct {
	PostID    string `json:"post_id"`
	AuthorDid string `json:"author_did"`
	Label     string `json:"label"`
}

type PostLike struct {
	PostID    string         `json:"post_id"`
	AuthorDid sql.NullString `json:"author_did"`
	LikeCount int64          `json:"like_count"`
}

type TopPoster struct {
	PostCount int64  `json:"post_count"`
	Handle    string `json:"handle"`
	AuthorDid string `json:"author_did"`
}
