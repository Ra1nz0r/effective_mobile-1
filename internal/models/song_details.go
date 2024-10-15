package models

// AddParams для получения данных и добавления песни.
type AddParams struct {
	Group string `json:"group,omitempty"`
	Song  string `json:"song,omitempty"`
}

// SongDetail для получения данных из внешнего API и обновления параметров песни.
type SongDetail struct {
	ID          int32  `json:"id,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}
