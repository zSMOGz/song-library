package utils

import (
	"path/filepath"
	"runtime"

	"song-library/internal/constants"
)

// GetProjectRoot возвращает абсолютный путь к корню проекта
func GetProjectRoot(callerSkip int) string {
	_, filename, _, _ := runtime.Caller(callerSkip)
	return filepath.Join(filepath.Dir(filename), constants.ProjectRootPath)
}
