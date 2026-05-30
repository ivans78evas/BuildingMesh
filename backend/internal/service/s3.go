package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type S3Service interface {
	Upload(reader io.Reader, fileName string) (string, error)
}

// LocalFileStorage implements S3Service for local development,
// ensuring heavy binary blobs are stored outside the SQLite database.
type LocalFileStorage struct {
	basePath string
}

func NewLocalFileStorage(path string) *LocalFileStorage {
	os.MkdirAll(path, 0755)
	return &LocalFileStorage{basePath: path}
}

func (s *LocalFileStorage) Upload(reader io.Reader, fileName string) (string, error) {
	fullPath := filepath.Join(s.basePath, fileName)

	out, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, reader)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("storage://%s", fileName), nil
}
