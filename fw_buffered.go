package alogw

import (
	"bufio"
	"os"
)

type fwBuffered struct {
	filename string
	bufsize  int
	f        *os.File
	bf       *bufio.Writer
	n        int64 // total bytes written
}

func (w fwBuffered) Filename() string {
	return w.filename
}

func (w *fwBuffered) Init(filename string) error {
	// f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	w.filename = filename
	f, err := createFile(filename)
	if err != nil {
		return err
	}
	w.f = f
	w.bf = bufio.NewWriterSize(f, w.bufsize)
	return nil
}

func (w *fwBuffered) Size() int64 {
	return w.n
}

func (w *fwBuffered) Write(p []byte) (int, error) {
	w.n += int64(len(p))
	return w.bf.Write(p)
}

func (w *fwBuffered) Close() error {
	w.bf.Flush()
	w.n = 0
	return w.f.Close()
}
