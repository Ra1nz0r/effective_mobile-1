package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ra1nz0r/effective_mobile-1/internal/logger"
)

// Добавляет ошибки в JSON и возвращает ответ в формате {"error":"ваш текст для ошибки"}.
func ErrReturn(err error, code int, w http.ResponseWriter) {
	result := make(map[string]string)
	result["error"] = err.Error()
	jsonResp, errJSON := json.Marshal(result)
	if errJSON != nil {
		//http.Error(w, errJSON.Error(), http.StatusInternalServerError)
		logger.Zap.Error(fmt.Errorf("failed attempt json-marshal response: %w", errJSON))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(code)

	if _, errWrite := w.Write(jsonResp); errWrite != nil {
		logger.Zap.Error("failed attempt WRITE response")
		return
	}
}
