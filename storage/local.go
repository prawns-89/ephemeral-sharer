package storage

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

// LocalStore manages files on the local disk
type LocalStore struct {
	BaseDir string
}

// New creates a LocalStore and ensures the folder exists
func New(dir string) *LocalStore {
	os.MkdirAll(dir, os.ModePerm)
	return &LocalStore{BaseDir: dir}
}

// Save takes a name and a stream of bytes, and writes it to disk
func (s *LocalStore) Save(filename string, content io.Reader) error {
	path := filepath.Join(s.BaseDir, filename)
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, content)
	return err
}

// Prune deletes files older than the given age
func (s *LocalStore) Prune(maxAge time.Duration) {
	files, _ := os.ReadDir(s.BaseDir)
	
	for _, file := range files {
		info, _ := file.Info()
		if time.Since(info.ModTime()) > maxAge {
			os.Remove(filepath.Join(s.BaseDir, file.Name()))
		}
	}
}