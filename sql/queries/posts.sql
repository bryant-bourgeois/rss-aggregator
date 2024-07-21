-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUserId :many
SELECT p.id, p.created_at, p.updated_at, p.title, p.url, p.description, p.published_at, p.feed_id FROM users u
LEFT JOIN feeds f ON u.id = f.user_id
LEFT JOIN posts p ON f.id = p.feed_id
WHERE u.id = $1
ORDER BY p.published_at
OFFSET $2 LIMIT $3;
