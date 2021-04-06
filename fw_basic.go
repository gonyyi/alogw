package alogw

import "os"

type fw struct {
	f        *os.File
	n        int64 // total bytes written
	filename string
}

func (w fw) Filename() string {
	return w.filename
}

func (w *fw) Init(filename string) error {
	// f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	w.filename = filename
	f, err := createFile(filename)
	if err != nil {
		return err
	}
	w.f = f
	return nil
}

func (w *fw) Size() int64 {
	return w.n
}

func (w *fw) Write(p []byte) (int, error) {
	w.n += int64(len(p))
	return w.f.Write(p)
}

func (w *fw) Close() error {
	w.n = 0
	return w.f.Close()
}
