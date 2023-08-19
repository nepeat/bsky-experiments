-- name: UpsertActor :exec
INSERT INTO actors (
        did,
        handle,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $3) ON CONFLICT (did) DO
UPDATE
SET handle = EXCLUDED.handle,
    updated_at = EXCLUDED.updated_at
WHERE actors.did = EXCLUDED.did;
-- name: GetActorByDID :one
SELECT *
FROM actors
WHERE did = $1;
-- name: GetActorByHandle :one
SELECT *
FROM actors
WHERE handle = $1;
-- name: FindActorsByHandle :many
SELECT *
FROM actors
WHERE handle ILIKE concat('%', $1, '%');
-- name: GetActorTypeAhead :many
SELECT did,
    handle,
    actors.created_at,
    CASE
        WHEN f.actor_did IS NOT NULL
        AND f2.actor_did IS NOT NULL THEN similarity(
            handle,
            concat('%', sqlc.arg('query')::text, '%')
        ) + 1
        WHEN f.actor_did IS NOT NULL THEN similarity(
            handle,
            concat('%', sqlc.arg('query')::text, '%')
        ) + 0.5
        ELSE similarity(
            handle,
            concat('%', sqlc.arg('query')::text, '%')
        )
    END::float AS score
FROM actors
    LEFT JOIN follows f ON f.target_did = did
    AND f.actor_did = sqlc.arg('actor_did')
    LEFT JOIN follows f2 ON f2.target_did = sqlc.arg('actor_did')
    AND f2.actor_did = did
WHERE handle ilike concat('%', sqlc.arg('query')::text, '%')
    AND NOT EXISTS (
        SELECT 1
        FROM blocks b
        WHERE (
                b.actor_did = sqlc.arg('actor_did')
                AND b.target_did = did
            )
            OR (
                b.actor_did = did
                AND b.target_did = sqlc.arg('actor_did')
            )
    )
ORDER BY score DESC
LIMIT $1;