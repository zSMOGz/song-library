// @title Song Library API
// @version 1.0
// @description API для работы с музыкальной библиотекой
// @host localhost:8080
// @BasePath /api
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	_ "song-library/docs"
	"song-library/internal/app"
	"song-library/internal/config"
	"song-library/internal/constants"
	_ "song-library/internal/handlers"
	"song-library/internal/utils"

	"github.com/joho/godotenv"
)

func main() {
	// Добавляем логгер с временными метками
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	// Получаем путь к корню проекта
	projectRoot := utils.GetProjectRoot(0)
	if err := godotenv.Load(filepath.Join(projectRoot, constants.EnvFileName)); err != nil {
		logger.Printf(constants.LogError, constants.ErrLoadingConfig, err)
	}
	logger.Println(constants.LogConfigLoaded)

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Printf(constants.LogError, constants.ErrLoadingConfig, err)
		logger.Fatal(constants.ErrInvalidData)
	}

	// Создаем новый экземпляр приложения
	app, err := app.NewApp(cfg, logger)
	if err != nil {
		logger.Printf(constants.LogError, constants.ErrAppInit, err)
		logger.Fatal(constants.ErrInvalidData)
	}

	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Настраиваем graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-quit
		logger.Printf(constants.LogSignalReceived, sig)
		cancel()
	}()

	// Запускаем приложение
	if err := app.Run(ctx); err != nil {
		logger.Printf(constants.LogError, constants.ErrAppRuntime, err)
		if shutdownErr := app.Shutdown(context.Background()); shutdownErr != nil {
			logger.Printf(constants.LogError, constants.ErrAppShutdown, shutdownErr)
		}
		os.Exit(1)
	}

	// Добавляем корректное завершение работы после выхода из Run
	if err := app.Shutdown(context.Background()); err != nil {
		logger.Printf(constants.LogError, constants.ErrAppShutdown, err)
		os.Exit(1)
	}
}
