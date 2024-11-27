package routers

import (
	"net/http"
	"song-library/internal/constants"
	"song-library/internal/handlers"
	"song-library/internal/logger"
	"song-library/internal/middleware"

	_ "song-library/docs" // автоматически сгенерированная документация

	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(songHandler *handlers.SongHandler, verseHandler *handlers.VerseHandler) http.Handler {
	router := http.NewServeMux()

	// Добавляем маршруты
	router.HandleFunc(constants.APISongsPath, songHandler.GetSongs)
	router.HandleFunc(constants.APISongDelete, songHandler.DeleteSong)
	router.HandleFunc(constants.APISongUpdate, songHandler.UpdateSong)
	router.HandleFunc(constants.APISongCreate, songHandler.CreateSong)
	router.HandleFunc(constants.APISongInfo, songHandler.GetSongInfo)
	router.HandleFunc(constants.APIVersesPath, verseHandler.GetVerses)
	router.Handle(constants.MetricsPath, promhttp.Handler())

	// Swagger
	router.HandleFunc(constants.SwaggerPath, httpSwagger.Handler(
		httpSwagger.URL(constants.SwaggerDocPath),
	))

	// Пприменяем middleware
	logger := logger.NewLogger()
	handler := middleware.RequestLogger(logger)(
		middleware.MetricsMiddleware(router),
	)

	return handler
}
