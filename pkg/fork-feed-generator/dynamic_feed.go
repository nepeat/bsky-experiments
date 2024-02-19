package feedgenerator

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bluesky-social/indigo/api/bsky"
	appbsky "github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/bytedance/sonic"
	"github.com/ericvolp12/bsky-experiments/pkg/consumer/store"
	"github.com/ericvolp12/bsky-experiments/pkg/events"
	"github.com/ericvolp12/bsky-experiments/pkg/sharddb"
	"github.com/labstack/gommon/log"
	"go.opentelemetry.io/otel"
)

type DynamicFeed struct {
	BSky         *events.BSky
	FeedName     string
	FeedActorDID string
	FeedJSON     string

	FilterMeta DynamicFeedMeta

	ShardDB       *sharddb.ShardDB
	FirehoseStore *store.Store
}

type DynamicFeedFilter struct {
	Type     string   `json:"type"`
	Operator string   `json:"operator"`
	Value    []string `json:"value"`
	Invert   bool     `json:"invert"` // optional field to invert the filter
}

type DynamicFeedMeta struct {
	Input   string              `json:"input"`
	Filters []DynamicFeedFilter `json:"filters"`
}

type FeedActor struct {
	ActorDID string
	Handle   string
}

func NewDynamicFeed(feedName string, feedActorDID string, feedJSON string, bskyClient *events.BSky, shardDB *sharddb.ShardDB, store *store.Store) (*DynamicFeed, []string, error) {
	// deserialice feedJSON to DynamicFeedMeta
	var filterMeta DynamicFeedMeta
	sonic.UnmarshalString(feedJSON, &filterMeta)

	// return the feed
	return &DynamicFeed{
		BSky:          bskyClient,
		FeedName:      feedName,
		FeedActorDID:  feedActorDID,
		FeedJSON:      feedJSON,
		FilterMeta:    filterMeta,
		ShardDB:       shardDB,
		FirehoseStore: store,
	}, []string{feedName}, nil
}

type ErrInvalidCursor struct {
	error
}

func ParseCursor(cursor string) (time.Time, string, string, error) {
	var err error
	createdAt := time.Now()
	actorDID := ""
	rkey := syntax.NewTID(time.Now().UnixMicro(), 0).String()

	if cursor != "" {
		cursorParts := strings.Split(cursor, "|")
		if len(cursorParts) != 3 {
			return createdAt, actorDID, rkey, ErrInvalidCursor{fmt.Errorf("cursor is invalid (wrong number of parts)")}
		}
		createdAtStr := cursorParts[0]

		createdAtUnixMili := int64(0)
		createdAtUnixMili, err = strconv.ParseInt(createdAtStr, 10, 64)
		if err != nil {
			return createdAt, actorDID, rkey, ErrInvalidCursor{fmt.Errorf("cursor is invalid (failed to parse createdAt)")}
		}

		createdAt = time.UnixMilli(createdAtUnixMili)

		actorDID = cursorParts[1]
		rkey = cursorParts[2]
	}

	return createdAt, actorDID, rkey, nil
}

func AssembleCursor(createdAt time.Time, actorDID string, rkey string) string {
	return fmt.Sprintf("%d|%s|%s", createdAt.UnixMilli(), actorDID, rkey)
}

func (plf *DynamicFeed) GetBulkActorsFromDIDs(ctx context.Context, actors map[string]string) map[string]string {
	result := map[string]string{}

	// Get all raw DIDs from actors
	actorDIDs := []string{}
	for actorDID, _ := range actors {
		actorDIDs = append(actorDIDs, actorDID)
	}

	// Attempt to get all from the DB first.
	dbActors, err := plf.FirehoseStore.Queries.GetActorByDIDs(ctx, actorDIDs)
	if err != nil {
		log.Errorf("error getting actors from DB: %v", err)
	}
	if len(dbActors) > 0 {
		for _, dbActor := range dbActors {
			result[dbActor.Did] = dbActor.Handle
			delete(actors, dbActor.Did)
		}
	}

	// Backfill the remaining actors from the BSky API
	for actorDID, _ := range actors {
		log.Warnf("Doing fallback to BSky for actor=%s", actorDID)
		handle, err := plf.BSky.ResolveDID(ctx, actorDID)
		if err != nil {
			log.Errorf("error resolving DID: %v", err)
			continue
		}
		result[actorDID] = handle
	}

	return result
}

func (plf *DynamicFeed) GetPage(ctx context.Context, feed string, userDID string, limit int64, cursor string) ([]*appbsky.FeedDefs_SkeletonFeedPost, *string, error) {
	tracer := otel.Tracer("DynamicFeed")
	ctx, span := tracer.Start(ctx, "GetPage")
	defer span.End()

	var err error
	createdAt := time.Now()
	var authorDID string
	var rkey string

	filteredPostURIs := []string{}
	newRkey := ""
	hasMore := false

	createdAt, authorDID, rkey, err = ParseCursor(cursor)

	// Otherwise use the sharddb to get the posts (high recent hit-rate)
	bucket, err := sharddb.GetBucketFromRKey(rkey)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting bucket from rkey: %w", err)
	}

	maxPages := 25
	pageSize := 15000
	hasMore = true

	metaPageCursor := createdAt
	for i := 0; i < maxPages; i++ {
		log.Infof("got new page for %s, page=%d, size=%d", plf.FeedName, i, len(filteredPostURIs))
		// Walk posts in reverse chronological order
		dbPosts, nextCursor, err := plf.ShardDB.GetPosts(ctx, bucket, pageSize, metaPageCursor)
		if err != nil {
			return nil, nil, fmt.Errorf("error getting post metas: %w", err)
		}

		// Pre-fetch the DID -> Actor mappings
		dbActors := map[string]string{}
		for _, post := range dbPosts {
			dbActors[post.ActorDID] = ""
		}
		feedActors := plf.GetBulkActorsFromDIDs(ctx, dbActors)

		// Pick out the posts from the non-moots
		for _, post := range dbPosts {
			createdAt = post.IndexedAt
			authorDID = post.ActorDID
			newRkey = post.Rkey
			postHandle := feedActors[post.ActorDID]

			if plf.FilterPost(ctx, post, postHandle) {
				filteredPostURIs = append(filteredPostURIs, fmt.Sprintf("at://%s/app.bsky.feed.post/%s", post.ActorDID, post.Rkey))
			}

			if len(filteredPostURIs) >= int(limit) {
				break
			}
		}

		if len(filteredPostURIs) >= int(limit) {
			break
		}

		metaPageCursor = nextCursor
		if nextCursor.IsZero() {
			bucket = bucket - 1
			metaPageCursor = time.Now()
			if bucket < 0 {
				hasMore = false
				break
			}
		}
	}

	log.Infof("finalized page, size=%d", len(filteredPostURIs))

	if newRkey == "" {
		// Set the rkey cursor to the next bucket
		newRkey = sharddb.GetHighestRKeyForBucket(bucket - 1)
	}

	// Convert to appbsky.FeedDefs_SkeletonFeedPost
	posts := []*appbsky.FeedDefs_SkeletonFeedPost{}
	for _, uri := range filteredPostURIs {
		posts = append(posts, &appbsky.FeedDefs_SkeletonFeedPost{Post: uri})
	}

	var next *string
	if hasMore {
		newCursor := AssembleCursor(createdAt, authorDID, newRkey)
		next = &newCursor
	}

	return posts, next, nil
}

func (plf *DynamicFeed) Describe(ctx context.Context) ([]appbsky.FeedDescribeFeedGenerator_Feed, error) {
	tracer := otel.Tracer("dynamic-feed")
	ctx, span := tracer.Start(ctx, "Describe")
	defer span.End()

	feeds := []appbsky.FeedDescribeFeedGenerator_Feed{
		{
			Uri: "at://" + plf.FeedActorDID + "/app.bsky.feed.generator/" + plf.FeedName,
		},
	}

	return feeds, nil
}

func (plf *DynamicFeed) FilterPost(ctx context.Context, post *sharddb.Post, handle string) bool {
	results := []bool{}

	for _, feedFilter := range plf.FilterMeta.Filters {
		var searchValue string
		result := false

		// get the specfic value from the post as needed
		if feedFilter.Type == "text" {
			postObject := bsky.FeedPost{}
			err := sonic.Unmarshal(post.Raw, &postObject)

			if err != nil {
				postString := string(post.Raw)
				// TODO: remove this debugging shit code
				fmt.Println(postString)
				fmt.Printf("error decoding json value: %v\n", err)
				fmt.Println("-------------------")
				return false
			}

			searchValue = postObject.Text
			searchValue = strings.ToLower(searchValue)
		} else if feedFilter.Type == "domain" {
			if handle == "" {
				results = append(results, false)
				continue
			}
			searchValue = handle
		}

		// search for the specfic data
		if feedFilter.Operator == "contains" {
			for _, value := range feedFilter.Value {
				if strings.Contains(searchValue, value) {
					result = true
					break
				}
			}
		} else if feedFilter.Operator == "endswith" {
			for _, value := range feedFilter.Value {
				if strings.HasSuffix(searchValue, value) {
					result = true
					break
				}
			}
		} else if feedFilter.Operator == "startswith" {
			for _, value := range feedFilter.Value {
				if strings.HasPrefix(searchValue, value) {
					result = true
					break
				}
			}
		}

		// invert if needed
		if feedFilter.Invert {
			result = !result
		}

		results = append(results, result)
	}

	// if any of the results are false, then the post is not valid
	for _, result := range results {
		if !result {
			return false
		}
	}

	return true
}
