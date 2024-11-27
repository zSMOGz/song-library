package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"song-library/internal/config"
	"song-library/internal/constants"
	"song-library/internal/db"
	"song-library/internal/migrations"
	"song-library/internal/server"
)

type App struct {
	cfg    *config.Config
	db     *db.Database
	server *http.Server
	logger *log.Logger
}

func NewApp(cfg *config.Config, logger *log.Logger) (*App, error) {
	database, err := db.NewDatabase(cfg.GetDBConnString())
	if err != nil {
		err = fmt.Errorf(constants.ErrDBConnection, err)
		logger.Println(err)
		return nil, err
	}

	return &App{
		cfg:    cfg,
		db:     database,
		logger: logger,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	// Инициализация и запуск миграций
	migrator, err := migrations.NewMigrator(ctx, a.cfg.DB, a.logger)
	if err != nil {
		return fmt.Errorf(constants.ErrMigratorInit+constants.ErrFormatAddition, err)
	}
	defer migrator.Close()

	if err := migrator.Up(ctx); err != nil {
		if downErr := migrator.Down(ctx); downErr != nil {
			a.logger.Printf(constants.LogMigrationFailure, constants.ErrMigrationDown, downErr)
		}
		return fmt.Errorf(constants.ErrMigrationUp+constants.ErrFormatAddition, err)
	}

	// Используем существующую настройку сервера
	server, err := server.Setup(a.cfg, a.logger)
	if err != nil {
		return fmt.Errorf(constants.ErrServerSetup+constants.ErrFormatAddition, err)
	}
	a.server = server

	// Запускаем HTTP сервер в горутине
	go func() {
		a.logger.Printf(constants.LogServerStarted, a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Printf(constants.LogError, constants.ErrServerCritical, err)
		}
	}()

	// Ожидаем сигнал завершения из контекста
	<-ctx.Done()
	a.logger.Println(constants.LogShutdownNotice)
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			a.logger.Printf(constants.LogError, constants.ErrGracefulShutdown, err)
			return fmt.Errorf(constants.ErrGracefulShutdown+constants.ErrFormatAddition, err)
		}
	}

	if a.db != nil {
		if err := a.db.Close(); err != nil {
			a.logger.Printf(constants.LogError, constants.ErrDBConnection, err)
			return fmt.Errorf(constants.ErrDBConnection, err)
		}
	}

	a.logger.Println(constants.LogServerStopped)
	return nil
}
