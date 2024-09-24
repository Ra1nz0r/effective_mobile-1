package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ra1nz0r/effective_mobile-1/internal/logger"
)

// Собирает все метрики метрики из локального хранилища и выводит их в
// результирующей карте при получении GET запроса.
// Вызывает метод интерфейса, который возвращает копию локального хранилища.
// Формат JSON, в виде {"Alloc":146464,"Frees":10,...}.
func GetAllMetrics(w http.ResponseWriter, r *http.Request) {

	ans, errJSON := json.Marshal("Hello from Web!")
	if errJSON != nil {
		logger.Zap.Error(fmt.Errorf("failed attempt json-marshal response: %w", errJSON))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if _, errWrite := w.Write([]byte(ans)); errWrite != nil {
		logger.Zap.Error("failed attempt WRITE response")
		return
	}
}
