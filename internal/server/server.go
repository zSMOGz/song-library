package server

import (
	"fmt"
	"log"
	"net/http"
	"song-library/internal/config"
	"song-library/internal/constants"
	"song-library/internal/db"
	"song-library/internal/handlers"
	"song-library/internal/repository"
	"song-library/internal/routers"
	"strings"
	"time"
)

func Setup(cfg *config.Config, logger *log.Logger) (*http.Server, error) {
	database, err := db.NewDatabase(cfg.GetDBConnString())
	if err != nil {
		logger.Printf(constants.LogError, constants.ErrDBConnection, err)
		return nil, fmt.Errorf(constants.ErrDBConnection, err)
	}

	songRepo, err := repository.NewSongRepository(database)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrFormat, constants.ErrSongRepoCreate, err)
	}

	verseRepo, err := repository.NewVerseRepository(database)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrFormat, constants.ErrVerseRepoCreate, err)
	}

	logger.Println(constants.LogReposInitialized)

	songHandler := handlers.NewSongHandler(songRepo, logger, cfg.ServerAddress)
	verseHandler := handlers.NewVerseHandler(verseRepo, logger)

	serverAddress := cfg.ServerAddress
	if idx := strings.Index(serverAddress, "//"); idx != -1 {
		serverAddress = serverAddress[idx+2:]
	}
	logger.Printf(constants.LogServerSetupAddr, serverAddress)

	return &http.Server{
		Addr:         serverAddress,
		Handler:      routers.SetupRoutes(songHandler, verseHandler),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}, nil
}
