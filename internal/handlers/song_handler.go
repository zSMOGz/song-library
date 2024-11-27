package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"errors"
	"log"
	"song-library/internal/constants"
	"song-library/internal/models"
	"song-library/internal/repository"
)

type SongHandler struct {
	repo          *repository.SongRepository
	logger        *log.Logger
	ServerAddress string
}

func NewSongHandler(repo *repository.SongRepository, logger *log.Logger, ServerAddress string) *SongHandler {
	return &SongHandler{
		repo:          repo,
		logger:        logger,
		ServerAddress: ServerAddress,
	}
}

// @Summary Получить список песен
// @Description Получить список песен с возможностью фильтрации
// @Tags songs
// @Accept json
// @Produce json
// @Param title query string false "Название песни"
// @Param artist query string false "Исполнитель"
// @Param album query string false "Альбом"
// @Param genre query string false "Жанр"
// @Param year query int false "Год выпуска"
// @Param page query int false "Номер страницы" default(1)
// @Param per_page query int false "Количество элементов на странице" default(10) maximum(100)
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {string} string "Ошибка валидации"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /songs [get]
func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Printf(constants.LogMethodNotAllowed, r.Method, constants.HandlerGetSongs)
		http.Error(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	filter := parseFilter(r)
	if err := validateFilter(filter); err != nil {
		h.logger.Printf(constants.LogValidationError, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.repo.ListSongs(filter)
	if err != nil {
		h.logger.Printf(constants.LogError, constants.ErrGettingSongs, err)
		http.Error(w, constants.ErrGettingSongs, http.StatusInternalServerError)
		return
	}

	w.Header().Set(constants.HeaderCacheControl, constants.CacheControlValue)
	w.Header().Set(constants.HeaderContentType, constants.HeaderContentTypeJSON)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf(constants.LogEncodingError, err)
		http.Error(w, constants.ErrEncodingResponse, http.StatusInternalServerError)
		return
	}
}

func parseFilter(r *http.Request) models.SongFilter {
	filter := models.SongFilter{
		Title:   r.URL.Query().Get(constants.QueryParamTitle),
		Artist:  r.URL.Query().Get(constants.QueryParamArtist),
		Album:   r.URL.Query().Get(constants.QueryParamAlbum),
		Genre:   r.URL.Query().Get(constants.QueryParamGenre),
		Page:    constants.DefaultPage,
		PerPage: constants.DefaultPageSize,
	}

	if year := r.URL.Query().Get(constants.QueryParamYear); year != "" {
		filter.Year, _ = strconv.Atoi(year)
	}
	if page := r.URL.Query().Get(constants.QueryParamPage); page != "" {
		filter.Page, _ = strconv.Atoi(page)
	}
	if perPage := r.URL.Query().Get(constants.QueryParamPerPage); perPage != "" {
		filter.PerPage, _ = strconv.Atoi(perPage)
	}

	return filter
}

func validateFilter(f models.SongFilter) error {
	if f.Page < 1 {
		return errors.New(constants.ErrInvalidPage)
	}
	if f.PerPage < 1 || f.PerPage > 100 {
		return errors.New(constants.ErrInvalidPerPage)
	}
	return nil
}

// @Summary Удалить песню
// @Description Удалить песню по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id query int true "ID песни"
// @Success 204 "Песня успешно удалена"
// @Failure 400 {string} string "Некорректный ID"
// @Failure 404 {string} string "Песня не найдена"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /songs/delete [delete]
func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.logger.Printf(constants.LogMethodNotAllowed, r.Method, constants.HandlerDeleteSong)
		http.Error(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get(constants.QueryParamID))
	if err != nil {
		h.logger.Printf(constants.LogInvalidID, err)
		http.Error(w, constants.ErrInvalidID, http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteSong(id); err != nil {
		if err == sql.ErrNoRows {
			h.logger.Printf(constants.LogSongNotFound, id)
			http.Error(w, constants.ErrSongNotFound, http.StatusNotFound)
			return
		}
		h.logger.Printf(constants.LogError, constants.ErrDeletingSong, err)
		http.Error(w, constants.ErrDeletingSong, http.StatusInternalServerError)
		return
	}

	h.logger.Printf(constants.LogSuccessDelete, id)
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Обновить песню
// @Description Обновить информацию о песне
// @Tags songs
// @Accept json
// @Produce json
// @Param id query int true "ID песни"
// @Param song body models.SongUpdate true "Данные песни"
// @Success 200 "Песня успешно обновлена"
// @Failure 400 {string} string "Некорректные данные"
// @Failure 404 {string} string "Песня не найдена"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /songs/update [put]
func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.logger.Printf(constants.LogMethodNotAllowed, r.Method, constants.HandlerUpdateSong)
		http.Error(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get(constants.QueryParamID))
	if err != nil {
		h.logger.Printf(constants.LogInvalidID, err)
		http.Error(w, constants.ErrInvalidID, http.StatusBadRequest)
		return
	}

	// Проверяем Content-Type
	contentType := r.Header.Get(constants.HeaderContentType)
	if contentType != constants.HeaderContentTypeJSON {
		h.logger.Printf(constants.LogInvalidContentType, contentType)
		http.Error(w, constants.ErrInvalidContentType, http.StatusBadRequest)
		return
	}

	// Ограничиваем размер тела запроса
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB limit

	var songUpdate models.SongUpdate
	if err := json.NewDecoder(r.Body).Decode(&songUpdate); err != nil {
		h.logger.Printf(constants.LogDecodingError, err)
		http.Error(w, constants.ErrInvalidData, http.StatusBadRequest)
		return
	}

	if songUpdate.Title == "" || songUpdate.Artist == "" {
		h.logger.Print(constants.LogMissingFields)
		http.Error(w, constants.ErrRequiredFields, http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateSong(id, songUpdate); err != nil {
		h.logger.Printf(constants.LogError, constants.ErrUpdatingSong, err)
		if err.Error() == constants.ErrSongNotFound {
			http.Error(w, constants.ErrSongNotFound, http.StatusNotFound)
			return
		}
		http.Error(w, constants.ErrUpdatingSong, http.StatusInternalServerError)
		return
	}

	h.logger.Printf(constants.LogSuccessUpdate, id)
	w.WriteHeader(http.StatusOK)
}

// @Summary Создать песню
// @Description Создать новую песню
// @Tags songs
// @Accept json
// @Produce json
// @Param input body models.SimpleSongInput true "Данные песни"
// @Success 201 {object} map[string]int "ID созданной песни"
// @Failure 400 {string} string "Некорректные данные"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /songs/create [post]
func (h *SongHandler) CreateSong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Printf(constants.LogMethodNotAllowed, r.Method, constants.HandlerCreateSong)
		http.Error(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	// Декодируем входящий JSON
	var input struct {
		Group string `json:"group"`
		Song  string `json:"song"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Printf(constants.LogDecodingError, err)
		http.Error(w, constants.ErrDecodingJSON, http.StatusBadRequest)
		return
	}

	// Запрос к API для получения дополнительной информации
	apiURL := fmt.Sprintf(constants.DefaultFormat,
		h.ServerAddress,
		fmt.Sprintf(constants.APIInfoURLFormat,
			url.QueryEscape(input.Group),
			url.QueryEscape(input.Song)))

	resp, err := http.Get(apiURL)
	if err != nil {
		h.logger.Printf(constants.LogError, constants.ErrFetchingSongInfo, err)
		http.Error(w, constants.ErrFetchingSongInfo, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var song models.Song
	if err := json.NewDecoder(resp.Body).Decode(&song); err != nil {
		h.logger.Printf(constants.LogError, constants.ErrProcessingSongInfo, err)
		http.Error(w, constants.ErrProcessingSongInfo, http.StatusInternalServerError)
		return
	}

	// Заполняем базовую информацию
	song.Title = input.Song
	song.Artist = input.Group

	// Сохраняем в базу данных
	id, err := h.repo.CreateSong(&song)
	if err != nil {
		h.logger.Printf(constants.LogError, constants.ErrSavingSong, err)
		http.Error(w, constants.ErrSavingSong, http.StatusInternalServerError)
		return
	}

	// Возвращаем ID созданной песни
	w.Header().Set(constants.HeaderCacheControl, constants.CacheControlValue)
	w.Header().Set(constants.HeaderContentType, constants.HeaderContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// @Summary Получить информацию о песне
// @Description Получить детальную информацию о песне по исполнителю и названию
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string true "Исполнитель"
// @Param song query string true "Название песни"
// @Success 200 {object} models.Song
// @Failure 400 {string} string "Некорректные параметры запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /songs/info [get]
func (h *SongHandler) GetSongInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Printf(constants.LogMethodNotAllowed, r.Method, constants.HandlerGetSongInfo)
		http.Error(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	filter := models.SongFilter{
		Title:   r.URL.Query().Get(constants.QueryParamSong),
		Artist:  r.URL.Query().Get(constants.QueryParamGroup),
		Page:    constants.DefaultPage,
		PerPage: constants.DefaultPageSize,
	}

	if err := validateFilter(filter); err != nil {
		h.logger.Printf(constants.LogValidationError, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.repo.ListSongs(filter)
	if err != nil {
		h.logger.Printf(constants.LogError, constants.ErrGettingSongs, err)
		http.Error(w, constants.ErrGettingSongs, http.StatusInternalServerError)
		return
	}

	w.Header().Set(constants.HeaderCacheControl, constants.CacheControlValue)
	w.Header().Set(constants.HeaderContentType, constants.HeaderContentTypeJSON)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf(constants.LogEncodingError, err)
		http.Error(w, constants.ErrEncodingResponse, http.StatusInternalServerError)
		return
	}
}
