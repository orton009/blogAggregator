-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * from users where name = $1;

-- name: DeleteAllUsers :exec
DELETE from users;

-- name: ListUsers :many
SELECT * from users;

-- name: GetUserById :one
SELECT * from users where id = $1;
