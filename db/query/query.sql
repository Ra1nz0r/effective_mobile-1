-- name: Add :one
INSERT INTO library ("group", song)
VALUES ($1, $2)
RETURNING *;
-- name: Delete :exec
DELETE FROM library
WHERE id = $1;
-- name: Fetch :exec
UPDATE library
SET "releaseDate" = $2,
    text = $3,
    link = $4
WHERE id = $1;
-- name: GetOne :one
SELECT *
FROM library
WHERE id = $1
LIMIT 1;
-- name: GetText :one
SELECT text
FROM library
WHERE id = $1
LIMIT 1;
-- name: ListAll :many
SELECT *
FROM library
ORDER BY id;
-- name: ListWithFilters :many
SELECT *
FROM library
WHERE (
        "group" ILIKE '%' || $1 || '%'
        OR $1 IS NULL
    )
    AND (
        song ILIKE '%' || $2 || '%'
        OR $2 IS NULL
    )
    AND (
        "releaseDate" >= $3
        OR $3 IS NULL
    )
    AND (
        "text" ILIKE '%' || $4 || '%'
        OR $4 IS NULL
    )
ORDER BY id
LIMIT $5 OFFSET $6;
-- name: Update :exec
UPDATE library
SET "releaseDate" = $2,
    "text" = $3,
    link = $4
WHERE id = $1;