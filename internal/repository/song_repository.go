package repository

import (
	"database/sql"
	"fmt"

	"embed"
	"song-library/internal/constants"
	"song-library/internal/db"
	"song-library/internal/models"
)

//go:embed queries/songs/*.sql
var songQueries embed.FS

type SongRepository struct {
	BaseRepository
	db *db.Database
}

func NewSongRepository(db *db.Database) (*SongRepository, error) {
	queries, err := loadQueries(songQueries, constants.SongQueriesPath)
	if err != nil {
		return nil, err
	}

	return &SongRepository{
		BaseRepository: BaseRepository{queries: queries},
		db:             db,
	}, nil
}

func (r *SongRepository) GetSong(id int) (*models.Song, error) {
	row := r.db.QueryRow(r.queries[constants.QueryGet], id)
	song := &models.Song{}
	err := row.Scan(&song.ID, &song.Title, &song.Artist, &song.Album,
		&song.Genre, &song.Duration, &song.ReleaseDate,
		&song.Text, &song.Link)
	if err != nil {
		return nil, err
	}
	return song, nil
}

func (r *SongRepository) ListSongs(filter models.SongFilter) (*models.PaginatedResponse, error) {
	offset := (filter.Page - 1) * filter.PerPage

	rows, err := r.db.Query(r.queries[constants.QueryListSongs],
		filter.Title, filter.Artist, filter.Album,
		filter.Year, filter.Genre,
		filter.PerPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	var totalCount int

	for rows.Next() {
		var s models.Song
		var nullText, nullLink, nullAlbum, nullGenre sql.NullString
		var nullReleaseDate sql.NullTime

		err := rows.Scan(&s.ID, &s.Title, &s.Artist, &nullAlbum,
			&nullGenre, &s.Duration, &nullReleaseDate,
			&nullText, &nullLink, &s.CreatedAt, &s.UpdatedAt,
			&totalCount)
		if err != nil {
			return nil, err
		}

		s.Text = nullText.String
		s.Link = nullLink.String
		s.Album = nullAlbum.String
		s.Genre = nullGenre.String
		if nullReleaseDate.Valid {
			s.ReleaseDate = nullReleaseDate.Time.String()
		}

		songs = append(songs, s)
	}

	totalPages := (totalCount + filter.PerPage - 1) / filter.PerPage

	return &models.PaginatedResponse{
		Data:       songs,
		Total:      totalCount,
		Page:       filter.Page,
		PerPage:    filter.PerPage,
		TotalPages: totalPages,
	}, nil
}

func (r *SongRepository) DeleteSong(id int) error {
	result, err := r.db.Exec(r.queries[constants.QueryDeleteSong], id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SongRepository) UpdateSong(id int, songUpdate models.SongUpdate) error {
	err := r.db.QueryRow(
		r.queries[constants.QueryUpdateSong],
		songUpdate.Title,
		songUpdate.Artist,
		songUpdate.Album,
		songUpdate.ReleaseDate,
		songUpdate.Text,
		songUpdate.Link,
		songUpdate.Genre,
		songUpdate.Duration,
		id,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return fmt.Errorf(constants.ErrSongNotFound)
	}
	return err
}

func (r *SongRepository) CreateSimpleSong(input *models.SimpleSongInput) (int, error) {
	var id int
	err := r.db.QueryRow(r.queries[constants.QueryCreateSimpleSong],
		input.Song,
		input.Group,
	).Scan(&id)
	return id, err
}

func (r *SongRepository) CreateSong(song *models.Song) (int, error) {
	var id int
	err := r.db.QueryRow(
		r.queries[constants.QueryCreateSong],
		song.Title,
		song.Artist,
		song.Album,
		song.ReleaseDate,
		song.Text,
		song.Link,
		song.Genre,
		song.Duration,
	).Scan(&id)
	return id, err
}
