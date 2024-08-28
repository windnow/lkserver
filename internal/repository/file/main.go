package file

import (
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type FileRepo struct {
	repo string
}

func (fr *FileRepo) GetFile(fileId string) (io.Reader, string, error) {
	filePath := filepath.Join(fr.repo, fr.getFileName(fileId))

	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}

	mimeType := mime.TypeByExtension(filepath.Ext(filePath))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return file, mimeType, nil
}

func (fr *FileRepo) getFileName(fileId string) string {
	files, err := os.ReadDir(fr.repo)
	if err != nil {
		return fileId
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			if fileName == fileId {
				return file.Name()
			}
		}
	}

	return fileId
}
