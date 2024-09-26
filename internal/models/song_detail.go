package models

// Структура для получения данных из внешнего API и обновления параметров песни.
type SongDetail struct {
	ID          int32  `json:"id"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
