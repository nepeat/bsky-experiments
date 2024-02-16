# BSky Experiments

This is a fork of the upstream bsky-experiments to make this run on my group infrastructure. I have no idea what I am doing but I'm hoping to do my own data science using the upstream as a starting point and running on my infrastructure for hammering convienence.

Some changes made to this fork are:

- `vmagent/` - VictoriaMetrics scraping job for my dev box to push metrics to my VM cluster.
- bridge mode is used instead of compose specific networks out of sheer lazyness and the oddities of running Docker networks that are primarily IPv6.

## Personal Notes

- `consumer` job runs at all times to fetch the firehose
- `plc` + `backfill-pump` is a combination that has been handy to make `consumer` fetch the entire BSky network.

## Getting Started

For the folks that are brave enough to run this without help, here are my notes that got the code running on my infrastructure. Beware breakages since this is code meant for not your infrastructure nor my infrastructure.

1. At minimum, have a Scylla cluster and Postgres endpoint. I will not describe the process because even I do not know how I set it up.

My setup is 3x virtual machines running on fairly OK metals (192GB RAM, 2x 2TB SSDs in ZFS RAID1, 2x Xeon Gold 6140s). General specs at time of writing is 16GB RAM, 8 cores, 256GB per VM.

Patroni is used for SQL DB failover and Consul is used for service discovery for the active database primary. DB VMs have Scylla and Postgres colocated with each other. This is not a great idea, please do not try this at home.

Once Postgres and Scylla is up, here are the bare minimum environment variables needed to make this run:

```env
POSTGRES_URL=postgres://bsky:hunter2@primary.16-sea1.service.consul:5433/bsky
SHARD_DB_NODES=[2001:db8::f01]:19042,[2001:db8::f02]:19042,[2001:db8::f03]:19042
REDIS_ADDRESS=host.docker.internal:6379
```

These env vars can be placed in:

- .envrc (useful for debugging on the shell)
- .env
- .consumer.env (symlink recomended to .env)
- .plc.env (symlink recomended to .env)

2. Create the Scylla keyspace. I know 2 replication is weak but given that this is my environmenet, YOLO.

```sh
cqlsh -e "CREATE KEYSPACE bsky WITH replication = {'class': 'NetworkTopologyStrategy', 'datacenter1': '2'}"
```

3. Import the Scylla schema.

```sh
cqlsh -e "SOURCE './db/scylla/schema.cql'"
```

4. Import DB schemas.

```sh
psql $POSTGRES_URL < pkg/consumer/store/schema/schema.sql 
psql $POSTGRES_URL < pkg/consumer/store/schema/sentiment.sql 
psql $POSTGRES_URL < pkg/consumer/store/schema/repocleanup.sql 
psql $POSTGRES_URL < pkg/consumer/store/schema/jazbot.sql 
psql $POSTGRES_URL < pkg/consumer/store/schema/auth.sql 
```

5. Start Redis.

**NOTE**: Redis will be listening on all IPs, you the reader should firewall off Redis.

```sh
make redis-up
```

6. Start firehose consumer.

```sh
make consumer-up
```

# Original README

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
