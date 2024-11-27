package repository

import (
	"embed"
	"fmt"
	"strings"

	"song-library/internal/constants"
)

type BaseRepository struct {
	queries map[string]string
}

func loadQueries(fs embed.FS, path string) (map[string]string, error) {
	queries := make(map[string]string)
	files, err := fs.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrReadingDirectory, err)
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), constants.SQLExtension) {
			continue
		}

		filePath := path + "/" + f.Name()
		content, err := fs.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf(constants.ErrReadingFile, f.Name(), err)
		}
		name := strings.TrimSuffix(f.Name(), constants.SQLExtension)
		queries[name] = string(content)
	}
	return queries, nil
}
