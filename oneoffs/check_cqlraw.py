import json
import os

from cassandra.cluster import Cluster

cluster = Cluster(os.environ["SHARD_DB_NODES"].split(","))

session = cluster.connect(keyspace="bsky")

rows = session.execute("SELECT * from posts_revchron")
i = 0
for post_row in rows:
    post = post_row.raw.decode("utf-8")
    post = json.loads(post)
    i += 1
    if i % 10000 == 0:
        print(i, post_row.indexed_at)
