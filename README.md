# BSky Experiments

This is a fork of the upstream bsky-experiments to make this run on my group infrastructure. I have no idea what I am doing but I'm hoping to do my own data science using the upstream as a starting point and running on my infrastructure for hammering convienence.

Some changes made to this fork are:

- `vmagent/` - VictoriaMetrics scraping job for my dev box to push metrics to my VM cluster.
- bridge mode is used instead of compose specific networks out of sheer lazyness and the oddities of running Docker networks that are primarily IPv6.

## Personal Notes

- `consumer` job runs at all times to fetch the firehose
- `plc` + `backfill-pump` is a combination that has been handy to make `consumer` fetch the entire BSky network.

## Original README

This repo has contains some fun Go experiments interacting with BlueSky via the AT Protocol.

The Makefile contains some useful targets for building and running the experiments.

Check out the contents of the `cmd` folder for some Go binaries that are interesting.

Don't expect any help running these, they change often and are not built to work in every environment and require external infrastructure to operate properly (Postgres, Redis, ScyllaDB, etc.)

Some of the experiments include:

- `consumer` - A firehose consumer and backfiller that indexes the contents of the ATProto firehose for search and filtering by other services.
- `feedgen-go` - A feed generator backend that serves many different feeds from the datastores populated by the consumer and other services.
- `graphd` - An in-memory graph database that tracks follows between users to make lookups and intersections faster in my other services.
- `indexer` - A worker that grabs posts from the DB created by the consumer and processes them for language detection, sentiment analysis, and object detection of images by sending them to various python services from the `python` folder.
- `jazbot` - A bot that interacts with users and sources data from the DBs populated by other services.
- `ts-layout` - A TypeScript service that handles ForceAtlas2 layouts of threads for visualization in my [Thread Visualizer](https://bsky.jazco.dev/thread)
- `object-detection` - A Python service that downloads images and detects objects in them.
- `sentiment` - A Python service that runs sentiment analysis on the textual content of posts.
- `plc` - A shallow mirror of `plc.directory` that contains the "latest" values for DID -> Handle resolution.
- `search` - A catch-all Go HTTP API that tracks statistics of AT Proto and has some useful endpoints that the other services make use of.
