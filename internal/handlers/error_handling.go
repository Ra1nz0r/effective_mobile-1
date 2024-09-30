package handlers

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/Ra1nz0r/effective_mobile-1/internal/logger"
)

// ErrReturn добавляет ошибки в JSON и возвращает ответ в формате {"error":"ваш текст для ошибки"}.
func ErrReturn(err error, code int, w http.ResponseWriter) {
	result := map[string]string{
		"error": err.Error(),
	}

	resJSON, err := json.Marshal(result)
	if err != nil {
		logger.Zap.Error(fmt.Errorf("failed attempt json-marshal response: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(code)

	if _, err = w.Write(resJSON); err != nil {
		logger.Zap.Error(fmt.Errorf("failed attempt WRITE response: %w", err))
		return
	}
}
