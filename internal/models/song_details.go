package models

// Структура для получения данных из внешнего API и обновления параметров песни.
type SongDetail struct {
	ID          int32  `json:"id,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}

type FilterSongsRequest struct {
	GroupName   *string `json:"group,omitempty"`       // Фильтр по группе
	SongName    *string `json:"song,omitempty"`        // Фильтр по названию песни
	ReleaseDate *string `json:"releaseDate,omitempty"` // Фильтр по дате релиза
	Text        *string `json:"text,omitempty"`        // Фильтр по тексту песни
	Limit       int32   `json:"limit"`                 // Лимит для пагинации
	Offset      int32   `json:"offset"`                // Смещение для пагинации
}
