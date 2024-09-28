package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrReturn добавляет ошибки в JSON и возвращает ответ в формате {"error":"ваш текст для ошибки"}.
func ErrReturn(err error, code int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(code)

	json.NewEncoder(w).Encode(
		map[string]string{
			"error": err.Error(),
		},
	)
}
