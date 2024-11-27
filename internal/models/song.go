package models

import (
	"time"
)

type Song struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Artist      string    `json:"artist"`
	Album       string    `json:"album"`
	Genre       string    `json:"genre"`
	Duration    int       `json:"duration"`
	ReleaseDate string    `json:"releaseDate,omitempty"`
	Text        string    `json:"text,omitempty"`
	Link        string    `json:"link,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

type SongFilter struct {
	Title   string
	Artist  string
	Album   string
	Year    int
	Genre   string
	Page    int
	PerPage int
}

type SongUpdate struct {
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album,omitempty"`
	Genre       string `json:"genre,omitempty"`
	Duration    int    `json:"duration"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}

type PaginatedResponse struct {
	Data       []Song `json:"data"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	TotalPages int    `json:"total_pages"`
}

// Добавим новую структуру для упрощенного формата
type SimpleSongInput struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}
