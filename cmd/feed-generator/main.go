package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ericvolp12/bsky-experiments/pkg/auth"
	"github.com/ericvolp12/bsky-experiments/pkg/consumer/store"
	feedgenerator "github.com/ericvolp12/bsky-experiments/pkg/feed-generator"
	"github.com/ericvolp12/bsky-experiments/pkg/feed-generator/endpoints"
	"github.com/ericvolp12/bsky-experiments/pkg/feeds/authorlabel"
	"github.com/ericvolp12/bsky-experiments/pkg/feeds/bangers"
	"github.com/ericvolp12/bsky-experiments/pkg/feeds/cluster"
	"github.com/ericvolp12/bsky-experiments/pkg/feeds/firehose"
	followersexp "github.com/ericvolp12/bsky-experiments/pkg/feeds/followers-exp"
	"github.com/ericvolp12/bsky-experiments/pkg/feeds/hot"
	"github.com/ericvolp12/bsky-experiments/pkg/feeds/pins"
	"github.com/ericvolp12/bsky-experiments/pkg/feeds/postlabel"
	"github.com/ericvolp12/bsky-experiments/pkg/graphd/client"
	"github.com/ericvolp12/bsky-experiments/pkg/search"
	"github.com/ericvolp12/bsky-experiments/pkg/sharddb"
	"github.com/ericvolp12/bsky-experiments/pkg/tracing"
	ginprometheus "github.com/ericvolp12/go-gin-prometheus"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

func main() {
	app := cli.App{
		Name:    "feed-generator",
		Usage:   "bluesky feed generator",
		Version: "0.1.0",
	}

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Usage:   "enable debug logging",
			Value:   false,
			EnvVars: []string{"DEBUG"},
		},
		&cli.IntFlag{
			Name:    "port",
			Usage:   "port to serve metrics on",
			Value:   8080,
			EnvVars: []string{"PORT"},
		},
		&cli.StringFlag{
			Name:    "redis-address",
			Usage:   "redis address for storing progress",
			Value:   "localhost:6379",
			EnvVars: []string{"REDIS_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "redis-prefix",
			Usage:   "redis prefix for storing progress",
			Value:   "fg",
			EnvVars: []string{"REDIS_PREFIX"},
		},
		&cli.StringFlag{
			Name:     "registry-postgres-url",
			Usage:    "postgres url for the registry database",
			Value:    "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
			Required: true,
			EnvVars:  []string{"REGISTRY_POSTGRES_URL"},
		},
		&cli.StringFlag{
			Name:     "firehose-postgres-url",
			Usage:    "postgres url for the firehose database",
			Value:    "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
			Required: true,
			EnvVars:  []string{"FIREHOSE_POSTGRES_URL"},
		},
		&cli.StringFlag{
			Name:    "graph-json-url",
			Usage:   "URL to the exported graph JSON",
			Value:   "https://s3.jazco.io/exported_graph_enriched.json",
			EnvVars: []string{"GRAPH_JSON_URL"},
		},
		&cli.StringFlag{
			Name:     "service-endpoint",
			Usage:    "URL that the feed generator will be available at",
			Required: true,
			EnvVars:  []string{"SERVICE_ENDPOINT"},
		},
		&cli.StringFlag{
			Name:     "feed-actor-did",
			Usage:    "DID of the feed actor",
			Required: true,
			EnvVars:  []string{"FEED_ACTOR_DID"},
		},
		&cli.StringFlag{
			Name:    "graphd-root",
			Usage:   "root of the graphd service",
			Value:   "http://localhost:1323",
			EnvVars: []string{"GRAPHD_ROOT"},
		},
		&cli.StringSliceFlag{
			Name:    "shard-db-nodes",
			Usage:   "list of scylla nodes for shard db",
			EnvVars: []string{"SHARD_DB_NODES"},
		},
	}

	app.Action = FeedGenerator

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func FeedGenerator(cctx *cli.Context) error {
	ctx := cctx.Context
	var logger *zap.Logger

	if cctx.Bool("debug") {
		logger, _ = zap.NewDevelopment()
		logger.Info("Starting logger in DEBUG mode...")
	} else {
		logger, _ = zap.NewProduction()
		logger.Info("Starting logger in PRODUCTION mode...")
	}

	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("failed to sync logger on teardown: %+v", err.Error())
		}
	}()

	sugar := logger.Sugar()

	sugar.Info("Reading config from environment...")

	// Registers a tracer Provider globally if the exporter endpoint is set
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		log.Println("initializing tracer...")
		// Start tracer with 20% sampling rate
		shutdown, err := tracing.InstallExportPipeline(ctx, "BSky-Feed-Generator-Go", 0.2)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := shutdown(ctx); err != nil {
				log.Fatal(err)
			}
		}()
	}

	postRegistry, err := search.NewPostRegistry(cctx.String("registry-postgres-url"))
	if err != nil {
		log.Fatalf("Failed to create PostRegistry: %v", err)
	}
	defer postRegistry.Close()

	store, err := store.NewStore(cctx.String("firehose-postgres-url"))
	if err != nil {
		log.Fatalf("Failed to create Store: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: cctx.String("redis-address"),
	})

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		log.Fatalf("failed to instrument redis with tracing: %+v\n", err)
	}

	// Test the connection to redis
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %+v\n", err)
	}

	h := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	graphdClient := client.NewClient(cctx.String("graphd-root"), &h)

	var shardDBClient *sharddb.ShardDB
	if len(cctx.StringSlice("shard-db-nodes")) > 0 {
		shardDBClient, err = sharddb.NewShardDB(ctx, cctx.StringSlice("shard-db-nodes"), slog.Default())
		if err != nil {
			log.Fatalf("Failed to create ShardDB: %v", err)
		}
	}

	feedActorDID := cctx.String("feed-actor-did")

	// Set the acceptable DIDs for the feed generator to respond to
	// We'll default to the feedActorDID and the Service Endpoint as a did:web
	serviceURL, err := url.Parse(cctx.String("service-endpoint"))
	if err != nil {
		log.Fatal(fmt.Errorf("error parsing service endpoint: %w", err))
	}

	serviceWebDID := "did:web:" + serviceURL.Hostname()

	log.Printf("service DID Web: %s", serviceWebDID)

	acceptableDIDs := []string{feedActorDID, serviceWebDID}

	feedGenerator, err := feedgenerator.NewFeedGenerator(ctx, feedActorDID, serviceWebDID, acceptableDIDs, cctx.String("service-endpoint"))
	if err != nil {
		log.Fatalf("Failed to create FeedGenerator: %v", err)
	}

	endpoints, err := endpoints.NewEndpoints(feedGenerator, cctx.String("graph-json-url"), postRegistry)
	if err != nil {
		log.Fatalf("Failed to create Endpoints: %v", err)
	}

	// Create a cluster feed
	clustersFeed, clusterFeedAliases, err := cluster.NewClusterFeed(ctx, feedActorDID, postRegistry, store)
	if err != nil {
		log.Fatalf("Failed to create ClusterFeed: %v", err)
	}
	feedGenerator.AddFeed(clusterFeedAliases, clustersFeed)

	// Create a postlabel feed
	postLabelFeed, postLabelFeedAliases, err := postlabel.NewPostLabelFeed(ctx, feedActorDID, postRegistry)
	if err != nil {
		log.Fatalf("Failed to create PostLabelFeed: %v", err)
	}
	feedGenerator.AddFeed(postLabelFeedAliases, postLabelFeed)

	// Create an authorlabel feed
	authorLabelFeed, authorLabelFeedAliases, err := authorlabel.NewAuthorLabelFeed(ctx, feedActorDID, postRegistry)
	if err != nil {
		log.Fatalf("Failed to create AuthorLabelFeed: %v", err)
	}
	feedGenerator.AddFeed(authorLabelFeedAliases, authorLabelFeed)

	// Create a firehose feed
	firehoseFeed, firehoseFeedAliases, err := firehose.NewFirehoseFeed(ctx, feedActorDID, postRegistry)
	if err != nil {
		log.Fatalf("Failed to create FirehoseFeed: %v", err)
	}
	feedGenerator.AddFeed(firehoseFeedAliases, firehoseFeed)

	// Create a Bangers feed
	bangersFeed, bangersFeedAliases, err := bangers.NewBangersFeed(ctx, feedActorDID, store, redisClient)
	if err != nil {
		log.Fatalf("Failed to create BangersFeed: %v", err)
	}
	feedGenerator.AddFeed(bangersFeedAliases, bangersFeed)

	// Create a What's Hot feed
	hotFeed, hotFeedAliases, err := hot.NewHotFeed(ctx, feedActorDID, store, redisClient)
	if err != nil {
		log.Fatalf("Failed to create HotFeed: %v", err)
	}
	feedGenerator.AddFeed(hotFeedAliases, hotFeed)

	// // Create a My Followers feed
	// followersFeed, followersFeedAliases, err := followers.NewFollowersFeed(ctx, feedActorDID, store, graphdClient, redisClient)
	// if err != nil {
	// 	log.Fatalf("Failed to create FollowersFeed: %v", err)
	// }
	// feedGenerator.AddFeed(followersFeedAliases, followersFeed)

	if shardDBClient != nil {
		// Create Experimental Followers feed
		followersExpFeed, followersExpFeedAliases, err := followersexp.NewFollowersFeed(ctx, feedActorDID, graphdClient, redisClient, shardDBClient, store)
		if err != nil {
			log.Fatalf("Failed to create FollowersExpFeed: %v", err)
		}
		feedGenerator.AddFeed(followersExpFeedAliases, followersExpFeed)
	}

	// Create a My Pins feed
	pinsFeed, pinsFeedAliases, err := pins.NewPinsFeed(ctx, feedActorDID, store)
	if err != nil {
		log.Fatalf("Failed to create PinsFeed: %v", err)
	}
	feedGenerator.AddFeed(pinsFeedAliases, pinsFeed)

	router := gin.New()

	router.Use(gin.Recovery())

	router.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			start := time.Now()
			// These can get consumed during request processing
			path := c.Request.URL.Path
			query := c.Request.URL.RawQuery
			c.Next()

			end := time.Now().UTC()
			latency := end.Sub(start)

			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				for _, e := range c.Errors.Errors() {
					logger.Error(e)
				}
			} else if path != "/metrics" {
				logger.Info(path,
					zap.Int("status", c.Writer.Status()),
					zap.String("method", c.Request.Method),
					zap.String("path", path),
					zap.String("query", query),
					zap.String("ip", c.ClientIP()),
					zap.String("user-agent", c.Request.UserAgent()),
					zap.String("time", end.Format(time.RFC3339)),
					zap.String("feedQuery", c.GetString("feedQuery")),
					zap.String("feedName", c.GetString("feedName")),
					zap.Int64("limit", c.GetInt64("limit")),
					zap.String("cursor", c.GetString("cursor")),
					zap.Duration("latency", latency),
					zap.String("user_did", c.GetString("user_did")),
				)
			}
		}
	}())

	router.Use(ginzap.RecoveryWithZap(logger, true))

	// Serve static files from the public folder
	router.Static("/public", "./public")
	router.Static("/assets", "./public/assets")

	// Plug in OTEL Middleware and skip metrics endpoint
	router.Use(
		otelgin.Middleware(
			"BSky-Feed-Generator-Go",
			otelgin.WithFilter(func(req *http.Request) bool {
				return req.URL.Path != "/metrics"
			}),
		),
	)

	// CORS middleware
	router.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"https://bsky.jazco.dev"},
			AllowMethods: []string{"GET", "OPTIONS"},
			AllowHeaders: []string{"Origin", "Content-Length", "Content-Type"},
			AllowOriginFunc: func(origin string) bool {
				u, err := url.Parse(origin)
				if err != nil {
					return false
				}
				// Allow localhost and localnet requests for localdev
				return u.Hostname() == "localhost" || u.Hostname() == "10.0.6.32"
			},
		},
	))

	// Prometheus middleware
	p := ginprometheus.NewPrometheus("gin", nil)
	p.Use(router)

	// Init a store provider for API keys
	storeProvider := auth.NewStoreProvider(store)

	auther, err := auth.NewAuth(
		500_000,
		time.Hour*12,
		"https://plc.directory",
		40,
		"did:web:feedsky.jazco.io",
		storeProvider,
	)
	if err != nil {
		log.Fatalf("Failed to create Auth: %v", err)
	}

	router.GET("/update_cluster_assignments", endpoints.UpdateClusterAssignments)
	router.GET("/.well-known/did.json", endpoints.GetWellKnownDID)

	// JWT Auth middleware
	router.Use(auther.AuthenticateGinRequestViaJWT)

	router.GET("/xrpc/app.bsky.feed.getFeedSkeleton", endpoints.GetFeedSkeleton)
	router.GET("/xrpc/app.bsky.feed.describeFeedGenerator", endpoints.DescribeFeedGenerator)

	// Create Admin routes
	adminRoutes := router.Group("/admin")
	{
		adminRoutes.GET("/feeds", endpoints.GetFeeds)
		adminRoutes.Static("/dashboard", "./public")
	}

	// API Key Auth Middleware
	router.Use(auther.AuthenticateGinRequestViaAPIKey)
	router.PUT("/assign_user_to_feed", endpoints.AssignUserToFeed)
	router.PUT("/unassign_user_from_feed", endpoints.UnassignUserFromFeed)
	router.GET("/feed_members", endpoints.GetFeedMembers)

	log.Printf("Starting server on port %d", cctx.Int("port"))
	return router.Run(fmt.Sprintf(":%d", cctx.Int("port")))
}
