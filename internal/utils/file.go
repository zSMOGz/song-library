package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"song-library/internal/constants"
)

func ReadQueryFile(dir, filename string) (string, error) {
	projectRoot := GetProjectRoot(0)
	fullPath := filepath.Join(projectRoot, dir, filename)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf(constants.ErrReadingQueryFile, filename, err)
	}
	return string(content), nil
}
