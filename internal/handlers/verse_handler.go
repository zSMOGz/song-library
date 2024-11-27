package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"log"
	"song-library/internal/constants"
	"song-library/internal/repository"
)

type VerseHandler struct {
	repo   *repository.VerseRepository
	logger *log.Logger
}

func NewVerseHandler(repo *repository.VerseRepository, logger *log.Logger) *VerseHandler {
	return &VerseHandler{repo: repo, logger: logger}
}

// @Summary Получить куплеты песни
// @Description Получить список куплетов для конкретной песни
// @Tags verses
// @Accept json
// @Produce json
// @Param song_id query int true "ID песни"
// @Param page query int false "Номер страницы" default(1)
// @Param page_size query int false "Количество элементов на странице" default(10) maximum(50)
// @Success 200 {array} models.Verse
// @Failure 400 {string} string "Некорректный ID песни"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /verses [get]
func (h *VerseHandler) GetVerses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Printf(constants.LogMethodNotAllowed, r.Method, constants.HandlerGetVerses)
		http.Error(w, constants.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	songID, err := strconv.Atoi(r.URL.Query().Get(constants.QueryParamSongID))
	if err != nil {
		h.logger.Printf(constants.LogInvalidID, r.URL.Query().Get(constants.QueryParamSongID))
		http.Error(w, constants.ErrInvalidID, http.StatusBadRequest)
		return
	}

	page := constants.DefaultPage
	if pageStr := r.URL.Query().Get(constants.QueryParamPage); pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = constants.DefaultPage
		}
	}

	pageSize := constants.DefaultPageSize
	if pageSizeStr := r.URL.Query().Get(constants.QueryParamPageSize); pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 || pageSize > constants.MaxPageSize {
			pageSize = constants.DefaultPageSize
		}
	}

	verses, err := h.repo.GetVerses(songID, page, pageSize)
	if err != nil {
		h.logger.Printf(constants.LogGettingVerses, songID, err)
		http.Error(w, constants.ErrGettingVerses, http.StatusInternalServerError)
		return
	}

	w.Header().Set(constants.HeaderContentType, constants.HeaderContentTypeJSON)
	json.NewEncoder(w).Encode(verses)
}
