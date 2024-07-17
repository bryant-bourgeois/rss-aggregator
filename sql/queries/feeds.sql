-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListFeeds :many
SELECT * FROM feeds
ORDER BY name;

-- name: FollowFeed :one
INSERT INTO users_feeds (id, user_id, feed_id, created_at, updated_at)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListFeedFollows :many
SELECT * FROM users_feeds WHERE user_id = $1;

-- name: GetFeedFollow :one
SELECT * FROM users_feeds WHERE id = $1;

-- name: UnfollowFeed :exec
DELETE FROM users_feeds WHERE id = $1;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY CASE WHEN last_fetched_at IS NULL 
THEN 0 ELSE 1 END, last_fetched_at 
LIMIT $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $2, last_fetched_at = $2
WHERE id = $1;
