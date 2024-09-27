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

// Функция для запроса к внешнему API
func FetchSongDetails(group, song string) (*models.SongDetail, error) {
	//apiURL := os.Getenv("EXTERNAL_API_URL") // TODOOOOOO
	apiURL := "http://localhost:7777/info"

	fullURL := fmt.Sprintf("%s?group=%s&song=%s", apiURL, url.PathEscape(group), url.PathEscape(song))

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

func RunMigrations(dbURL string) error {
	m, err := migrate.New(
		"file://db/migration", // путь до вашей папки с миграциями // вынести в конфиг или енв путь до миграций
		dbURL)
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

// Проверка существования таблицы
func TableExists(db *sql.DB, tableName string) (bool, error) {
	var exists bool
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT FROM pg_tables
			WHERE schemaname = 'public' AND tablename = '%s'
		);`, tableName)
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
