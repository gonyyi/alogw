package alogw

import (
	"os"
	"path/filepath"
)

type writer interface {
	Filename() string
	Init(filename string) error
	Write(p []byte) (int, error)
	Size() int64
	Close() error
}

func createFile(fpath string) (*os.File, error) {
	// Create directory
	dir := filepath.Dir(fpath)
	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0744); err != nil { // rwxr--r--
			return nil, ErrCannotCreateDir
		}
	}

	// Create bufwriter
	return os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // rw-r--r--
}
