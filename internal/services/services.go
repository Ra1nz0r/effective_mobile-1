package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Ra1nz0r/effective_mobile-1/internal/models"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// FetchSongDetails делает запрос во внешний API и возвращает полученные сведения.
// Формат запроса: http://localhost:7777/info?group='group name'&song='song name'.
// Возвращает ошибку в случае неудачи.
func FetchSongDetails(group, song, externalApiURL string) (*models.SongDetail, error) {
	fullURL := fmt.Sprintf("%s?group=%s&song=%s", externalApiURL, url.PathEscape(group), url.PathEscape(song))

	resp, errResp := http.Get(fullURL)
	if errResp != nil {
		return nil, errResp
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %s", resp.Status)
	}

	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		return nil, fmt.Errorf("%w", errBody)
	}

	var songDetail models.SongDetail
	if errJson := json.Unmarshal(body, &songDetail); errJson != nil {
		return nil, fmt.Errorf("%w", errJson)
	}

	return &songDetail, nil
}

// RunMigrations запускает миграцию Up по указанному пути.
func RunMigrations(databaseURL, migrationPath string) error {
	m, err := migrate.New(migrationPath, databaseURL)
	if err != nil {
		return err
	}

	// Применение миграций
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// TableExists проверяет существование table в базе данных.
func TableExists(db *sql.DB, tableName string) (bool, error) {
	var exists bool
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT FROM pg_tables
			WHERE schemaname = 'public' OR schemaname = 'private'
			AND tablename = '%s'
		);`, tableName)
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
