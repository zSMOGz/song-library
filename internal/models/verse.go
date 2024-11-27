package models

import "time"

type Verse struct {
	ID          int       `json:"id"`
	SongID      int       `json:"song_id"`
	VerseNumber int       `json:"verse_number"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}
