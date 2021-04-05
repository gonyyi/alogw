package alogw

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
)

type fileWriter interface {
	Filename() string
	Init(filename string) error
	Write(p []byte) (int, error)
	Size() int
	Close() error
}

type fw struct {
	f        *os.File
	n        int // total bytes written
	filename string
}

func (w fw) Filename() string {
	return w.filename
}
func (w *fw) Init(filename string) error {
	// f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	w.filename = filename
	f, err := createFileIfNotExist(filename)
	if err != nil {
		return err
	}
	w.f = f
	return nil
}

func (w *fw) Size() int {
	return w.n
}

func (w *fw) Write(p []byte) (int, error) {
	w.n += len(p)
	return w.f.Write(p)
}

func (w *fw) Close() error {
	w.n = 0
	return w.f.Close()
}

type bfw struct {
	filename string
	bufsize  int
	f        *os.File
	bf       *bufio.Writer
	n        int // total bytes written
}

func (w bfw) Filename() string {
	return w.filename
}

func (w *bfw) Init(filename string) error {
	// f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	w.filename = filename
	f, err := createFileIfNotExist(filename)
	if err != nil {
		return err
	}
	w.f = f
	w.bf = bufio.NewWriterSize(f, w.bufsize)
	return nil
}

func (w *bfw) Size() int {
	return w.n
}

func (w *bfw) Write(p []byte) (int, error) {
	w.n += len(p)
	return w.bf.Write(p)
}

func (w *bfw) Close() error {
	w.bf.Flush()
	w.n = 0
	return w.f.Close()
}

const (
	errCannotCreateDir_err = `cannot create directory: %s`
)

func createFileIfNotExist(fpath string) (*os.File, error) {
	// Create directory
	dir := filepath.Dir(fpath)
	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0744); err != nil { // rwxr--r--
			return nil, fmt.Errorf(errCannotCreateDir_err, err.Error())
		}
	}

	// Create bufwriter
	return os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // rw-r--r--
}

type gzfw struct {
	filename string
	bufsize  int
	f        *os.File
	gz       *gzip.Writer
	bw       *bufio.Writer
	n        int // total bytes written
}

func (w gzfw) Filename() string {
	return w.filename
}

func (w *gzfw) Init(filename string) error {
	// f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	w.filename = filename
	f, err := createFileIfNotExist(filename)
	if err != nil {
		return err
	}
	w.f = f
	w.gz = gzip.NewWriter(f)
	w.bw = bufio.NewWriterSize(w.gz, w.bufsize)
	return nil
}

func (w *gzfw) Size() int {
	return w.n
}

func (w *gzfw) Write(p []byte) (int, error) {
	w.n += len(p)
	return w.bw.Write(p)
}

func (w *gzfw) Close() (err error) {
	if err := w.bw.Flush(); err!=nil{
		err = fmt.Errorf("e1:%w", err)
		println(1, err.Error())
	}
	if err := w.gz.Flush(); err != nil {
		err = fmt.Errorf("e2:%w", err)
	}
	if err := w.gz.Close(); err != nil {
		err = fmt.Errorf("e3:%w", err)
	}
	w.n = 0
	return w.f.Close()
}
