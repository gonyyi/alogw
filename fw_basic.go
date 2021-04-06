package alogw

import "os"

type fwBasic struct {
	f        *os.File
	n        int64 // total bytes written
	filename string
}

func (w fwBasic) Filename() string {
	return w.filename
}

func (w *fwBasic) Init(filename string) error {
	// f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	w.filename = filename
	f, err := createFile(filename)
	if err != nil {
		return err
	}
	w.f = f
	return nil
}

func (w *fwBasic) Size() int64 {
	return w.n
}

func (w *fwBasic) Write(p []byte) (int, error) {
	w.n += int64(len(p))
	return w.f.Write(p)
}

func (w *fwBasic) Close() error {
	w.n = 0
	return w.f.Close()
}
