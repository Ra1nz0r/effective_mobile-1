package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ra1nz0r/effective_mobile-1/internal/models"
	"github.com/Ra1nz0r/effective_mobile-1/internal/services"
	"github.com/stretchr/testify/assert"
)

// MockSongDetail это структура для тестирования
var MockSongDetail = models.SongDetail{
	Text:        "Some additional info about the song",
	Link:        "http://example.com/more-info",
	ReleaseDate: "10.10.2006",
}

func TestFetchSongDetails_Success(t *testing.T) {
	// Создаем тестовый HTTP-сервер, который будет имитировать внешний API.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что правильные параметры передаются в запрос.
		assert.Equal(t, "Muse", r.URL.Query().Get("group"))
		assert.Equal(t, "Supermassive Black Hole", r.URL.Query().Get("song"))

		// Возвращаем успешный JSON-ответ
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(MockSongDetail)
	}))
	defer mockServer.Close()

	// Вызываем тестируемую функцию
	result, err := services.FetchSongDetails("Muse", "Supermassive Black Hole", mockServer.URL)

	// Проверяем, что ошибки нет
	assert.NoError(t, err)
	// Проверяем, что результат соответствует ожидаемым данным
	assert.Equal(t, MockSongDetail, *result)
}

func TestFetchSongDetails_HttpError(t *testing.T) {
	// Закрываем HTTP-сервер для имитации ошибки сети.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	mockServer.Close()

	// Вызываем тестируемую функцию
	_, err := services.FetchSongDetails("Muse", "Supermassive Black Hole", mockServer.URL)

	// Ожидаем ошибку запроса
	assert.Error(t, err)
}

func TestFetchSongDetails_StatusNotOK(t *testing.T) {
	// Создаем тестовый HTTP-сервер, который возвращает статус 500.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	// Вызываем тестируемую функцию
	_, err := services.FetchSongDetails("Muse", "Supermassive Black Hole", mockServer.URL)

	// Ожидаем ошибку из-за неправильного статуса
	assert.EqualError(t, err, "API returned status: 500 Internal Server Error")
}

func TestFetchSongDetails_BodyReadError(t *testing.T) {
	// Создаем тестовый HTTP-сервер, который возвращает некорректное тело ответа.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// Закрываем соединение перед тем, как тело будет считано, что вызовет ошибку.
		w.(http.Flusher).Flush()
		r.Context()
	}))
	defer mockServer.Close()

	// Вызываем тестируемую функцию
	_, err := services.FetchSongDetails("Muse", "Supermassive Black Hole", mockServer.URL)

	// Ожидаем ошибку при чтении тела ответа
	assert.Error(t, err)
}

func TestFetchSongDetails_JSONUnmarshalError(t *testing.T) {
	// Создаем тестовый HTTP-сервер, который возвращает некорректный JSON.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer mockServer.Close()

	// Вызываем тестируемую функцию
	_, err := services.FetchSongDetails("Muse", "Supermassive Black Hole", mockServer.URL)

	// Ожидаем ошибку при парсинге JSON
	assert.Error(t, err)
}
