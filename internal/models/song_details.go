package models

// Структура для получения данных из внешнего API и обновления параметров песни.
type SongDetail struct {
	ID          int32  `json:"id,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}
