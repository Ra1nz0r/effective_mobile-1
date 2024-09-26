package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Ra1nz0r/effective_mobile-1/internal/models"
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
