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
	"github.com/ericvolp12/bsky-experiments/pkg/events"
	feedgenerator "github.com/ericvolp12/bsky-experiments/pkg/fork-feed-generator"
	"github.com/ericvolp12/bsky-experiments/pkg/fork-feed-generator/endpoints"
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
			Name:     "firehose-postgres-url",
			Usage:    "postgres url for the firehose database",
			Value:    "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
			Required: true,
			EnvVars:  []string{"FIREHOSE_POSTGRES_URL"},
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
		shutdown, err := tracing.InstallExportPipeline(ctx, "BSky-Feed-Generator-Go", 0.01)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := shutdown(ctx); err != nil {
				log.Fatal(err)
			}
		}()
	}

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

	var shardDBClient *sharddb.ShardDB
	shardDBClient, err = sharddb.NewShardDB(ctx, cctx.StringSlice("shard-db-nodes"), slog.Default())
	if err != nil {
		log.Fatalf("Failed to create ShardDB: %v", err)
	}

	// Create BSky Client
	dbConnectionString := ""
	postRegistryEnabled := false
	workerCount := 2
	plcDirectoryMirror := "http://localhost:8097"
	// plcDirectoryMirror := "https://plc.jazco.io"

	bsky, err := events.NewBSky(
		ctx,
		true,
		postRegistryEnabled,
		dbConnectionString,
		"graph_builder",
		plcDirectoryMirror,
		redisClient,
		workerCount,
	)
	if err != nil {
		log.Fatal(err)
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

	feedGenerator, err := feedgenerator.NewFeedGenerator(ctx, serviceWebDID, acceptableDIDs, cctx.String("service-endpoint"))
	if err != nil {
		log.Fatalf("Failed to create FeedGenerator: %v", err)
	}

	endpoints, err := endpoints.NewEndpoints(feedGenerator)
	if err != nil {
		log.Fatalf("Failed to create Endpoints: %v", err)
	}

	// Register static feeds
	feedGenerator.RegisterStaticFeed(feedActorDID, "dot_com", `{"input":"12h","filters":[{"type":"domain","operator":"endswith","value":[".com"]}]}`, bsky, shardDBClient, store)
	feedGenerator.RegisterStaticFeed(feedActorDID, "dot_net", `{"input":"12h","filters":[{"type":"domain","operator":"endswith","value":[".net"]}]}`, bsky, shardDBClient, store)
	feedGenerator.RegisterStaticFeed(feedActorDID, "dot_org", `{"input":"12h","filters":[{"type":"domain","operator":"endswith","value":[".org"]}]}`, bsky, shardDBClient, store)
	feedGenerator.RegisterStaticFeed(feedActorDID, "dot_me", `{"input":"12h","filters":[{"type":"domain","operator":"endswith","value":[".me"]}]}`, bsky, shardDBClient, store)
	feedGenerator.RegisterStaticFeed(feedActorDID, "dot_io", `{"input":"12h","filters":[{"type":"domain","operator":"endswith","value":[".io"]}]}`, bsky, shardDBClient, store)
	// feedGenerator.RegisterStaticFeed(feedActorDID, "avoid5", `{"input":"12h","filters":[{"type":"text","operator":"contains","invert": false,"value":["e","E"]}]}`, bsky, shardDBClient, store)
	feedGenerator.RegisterStaticFeed(feedActorDID, "testfeed", `{"input":"12h","filters":[{"type":"text","operator":"contains","invert": false,"value":["shit", "piss", "fuck", "cunt", "cocksucker", "motherfucker", "tits"]}]}`, bsky, shardDBClient, store)
	feedGenerator.RegisterStaticFeed(feedActorDID, "novowels", `{"input":"12h","filters":[{"type":"text","operator":"contains","invert": true,"value":["a","e","i","o","u","A","E","I","O","U"]}]}`, bsky, shardDBClient, store)
	// feedGenerator.RegisterStaticFeed(feedActorDID, "govfeed", `{"input":"7d","filters":[{"type":"domain","operator":"endswith","value":[".gov",".edu"]}]}`, bsky, shardDBClient, store)

	// if shardDBClient != nil {
	// 	// Create Experimental Followers feed
	// 	followersExpFeed, followersExpFeedAliases, err := followersexp.NewFollowersFeed(ctx, feedActorDID, graphdClient, redisClient, shardDBClient, store)
	// 	if err != nil {
	// 		log.Fatalf("Failed to create FollowersExpFeed: %v", err)
	// 	}
	// 	feedGenerator.AddFeed(followersExpFeedAliases, followersExpFeed)
	// }

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
			// TODO: this is a placeholder
			AllowOrigins: []string{"https://bsky.owo.me"},
			AllowMethods: []string{"GET", "OPTIONS"},
			AllowHeaders: []string{"Origin", "Content-Length", "Content-Type"},
			AllowOriginFunc: func(origin string) bool {
				u, err := url.Parse(origin)
				if err != nil {
					return false
				}
				// Allow localhost and localnet requests for localdev
				return u.Hostname() == "localhost" || u.Hostname() == "10.3.2.13"
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
		"did:web:feed-api.owo.me",
		storeProvider,
	)
	if err != nil {
		log.Fatalf("Failed to create Auth: %v", err)
	}

	router.GET("/.well-known/did.json", endpoints.GetWellKnownDID)

	// JWT Auth middleware
	router.Use(auther.AuthenticateGinRequestViaJWT)

	router.GET("/xrpc/app.bsky.feed.getFeedSkeleton", endpoints.GetFeedSkeleton)
	router.GET("/xrpc/app.bsky.feed.describeFeedGenerator", endpoints.DescribeFeedGenerator)

	// API Key Auth Middleware
	router.Use(auther.AuthenticateGinRequestViaAPIKey)

	log.Printf("Starting server on port %d", cctx.Int("port"))
	return router.Run(fmt.Sprintf(":%d", cctx.Int("port")))
}
