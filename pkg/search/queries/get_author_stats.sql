-- name: GetAuthorStats :one
WITH postcounts AS (
    SELECT
        author_did,
        COUNT(id) AS num_posts
    FROM
        posts
    GROUP BY
        author_did
),
percentiles AS (
    SELECT
        (percentile_cont(0.25) WITHIN GROUP (ORDER BY num_posts) * 1)::bigint AS p25,
        (percentile_cont(0.50) WITHIN GROUP (ORDER BY num_posts) * 1)::bigint AS p50,
        (percentile_cont(0.75) WITHIN GROUP (ORDER BY num_posts) * 1)::bigint AS p75,
        (percentile_cont(0.90) WITHIN GROUP (ORDER BY num_posts) * 1)::bigint AS p90,
        (percentile_cont(0.95) WITHIN GROUP (ORDER BY num_posts) * 1)::bigint AS p95,
        (percentile_cont(0.99) WITHIN GROUP (ORDER BY num_posts) * 1)::bigint AS p99
    FROM
        postcounts
),
counts AS (
    SELECT
        count(num_posts) AS total_authors,
        SUM(CASE WHEN num_posts > 1 THEN 1 ELSE 0 END) AS gt_1,
        SUM(CASE WHEN num_posts > 5 THEN 1 ELSE 0 END) AS gt_5,
        SUM(CASE WHEN num_posts > 10 THEN 1 ELSE 0 END) AS gt_10,
        SUM(CASE WHEN num_posts > 20 THEN 1 ELSE 0 END) AS gt_20,
        SUM(CASE WHEN num_posts > 100 THEN 1 ELSE 0 END) AS gt_100,
        SUM(CASE WHEN num_posts > 1000 THEN 1 ELSE 0 END)as gt_1000
    FROM
        postcounts
),
hellthread_posts AS (
    SELECT
        count(id) AS hellthread_post_count
    FROM
        posts
    WHERE
        root_post_id = '3juzlwllznd24'
)
SELECT
    total_authors,
    (SELECT SUM(num_posts) FROM postcounts)::bigint AS total_posts,
    hellthread_post_count,
    (SELECT AVG(num_posts) FROM postcounts)::float AS mean_posts_per_author,
    gt_1,
    gt_5,
    gt_10,
    gt_20,
    gt_100,
    gt_1000,
    p25,
    p50,
    p75,
    p90,
    p95,
    p99
FROM
    counts,
    percentiles,
    hellthread_posts
LIMIT 1;
