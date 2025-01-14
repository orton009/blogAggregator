-- CreatePosts :one
INSERT INTO posts (feed_id, title, url, description, published_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- GetPosts :many
SELECT * FROM posts
WHERE feed_id = $1
ORDER BY published_at DESC
LIMIT $2; 