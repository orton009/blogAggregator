-- name: CreateFeed :one
INSERT INTO feeds(id, name, url, created_at, updated_at, user_id) VALUES (
  $1, 
  $2, $3, $4, $5, $6
) RETURNING *;

-- name: DeleteAllFeeds :exec
DELETE from feeds;

-- name: ListFeeds :many
SELECT * from feeds;

-- name: GetFeed :one
SELECT * from feeds where url = $1;

-- name: MarkFeedAsFetched :exec
UPDATE feeds SET last_fetched_at = NOW(), updated_at = NOW() WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * from feeds ORDER BY last_fetched_at ASC NULLS FIRST LIMIT 1 ;