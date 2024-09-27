package db

import (
	"context"
	"database/sql"
)

const listLibrary = `-- name: ListLibrary :many
SELECT id, "group", song, "releaseDate", text, link
FROM library
WHERE (
        LOWER("group") LIKE '%' || LOWER($1) || '%' OR $1 IS NULL
    )
    AND (
        LOWER(song) LIKE '%' || LOWER($2) || '%' OR $2 IS NULL
    )
    AND (
        "releaseDate" >= $3 OR $3 IS NULL
    )
    AND (
        LOWER(text) LIKE '%' || LOWER($4) || '%' OR $4 IS NULL
    )
ORDER BY id
LIMIT $5 OFFSET $6;`

type ListLibraryParams struct {
	Group       sql.NullString
	Song        sql.NullString
	ReleaseDate sql.NullTime
	Text        sql.NullString
	Limit       int32
	Offset      int32
}

func (q *Queries) ListLibrary(ctx context.Context, arg ListLibraryParams) ([]Library, error) {
	rows, err := q.db.QueryContext(ctx, listLibrary,
		arg.Group,
		arg.Song,
		arg.ReleaseDate,
		arg.Text,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Library
	for rows.Next() {
		var i Library
		if err := rows.Scan(
			&i.ID,
			&i.Group,
			&i.Song,
			&i.ReleaseDate,
			&i.Text,
			&i.Link,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
