-- name: AddArtist :one
INSERT INTO artist ("group")
VALUES ($1)
RETURNING *;
-- name: AddSongWithID :one
INSERT INTO library (group_id, "song")
VALUES ($1, $2)
RETURNING *;
-- name: CheckSongWithID :one
SELECT EXISTS (
        SELECT 1
        FROM library
        WHERE group_id = $1
            AND song = $2
    );
-- name: Delete :exec
DELETE FROM library
WHERE id = $1;
-- name: Fetch :exec
UPDATE library
SET "releaseDate" = $2,
    text = $3,
    link = $4
WHERE id = $1;
-- name: GetArtistID :one
SELECT id
FROM artist
WHERE "group" = $1
LIMIT 1;
-- name: GetOne :one
SELECT *
FROM library
WHERE id = $1
LIMIT 1;
-- name: GetText :one
SELECT library.id,
    artist."group",
    library.song,
    library.text
FROM library
    JOIN artist ON library.group_id = artist.id
WHERE library.id = $1
LIMIT 1;
-- name: ListWithFilters :many 
SELECT library.id,
    artist."group",
    library.song,
    library."releaseDate",
    library.text,
    library.link
FROM library
    JOIN artist ON library.group_id = artist.id
WHERE (
        artist."group" ILIKE '%' || $1 || '%'
        OR $1 IS NULL
    )
    AND (
        library.song ILIKE '%' || $2 || '%'
        OR $2 IS NULL
    )
    AND (
        library."releaseDate" >= $3
        OR $3 IS NULL
    )
    AND (
        library."text" ILIKE '%' || $4 || '%'
        OR $4 IS NULL
    )
ORDER BY library.id
LIMIT $5 OFFSET $6;
-- name: Update :exec
UPDATE library
SET "releaseDate" = COALESCE(
        NULLIF($2::date, '0001-01-01'::date),
        "releaseDate"
    ),
    "text" = COALESCE(NULLIF($3, ''), "text"),
    link = COALESCE(NULLIF($4, ''), link)
WHERE id = $1;