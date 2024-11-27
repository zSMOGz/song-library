package repository

import (
	"embed"
	"song-library/internal/constants"
	"song-library/internal/db"
	"song-library/internal/models"
)

//go:embed queries/verses/*.sql
var verseQueries embed.FS

type VerseRepository struct {
	BaseRepository
	db *db.Database
}

func NewVerseRepository(db *db.Database) (*VerseRepository, error) {
	queries, err := loadQueries(verseQueries, constants.VerseQueriesPath)
	if err != nil {
		return nil, err
	}

	return &VerseRepository{
		BaseRepository: BaseRepository{queries: queries},
		db:             db,
	}, nil
}

func (r *VerseRepository) GetVerses(songID int, page, pageSize int) ([]models.Verse, error) {
	offset := (page - 1) * pageSize

	rows, err := r.db.Query(r.queries[constants.QueryGet], songID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []models.Verse
	for rows.Next() {
		var v models.Verse
		if err := rows.Scan(&v.ID, &v.SongID, &v.VerseNumber, &v.Content, &v.CreatedAt); err != nil {
			return nil, err
		}
		verses = append(verses, v)
	}
	return verses, nil
}
