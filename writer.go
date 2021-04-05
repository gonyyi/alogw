package alogw

import (
	"sync"
)

// CONST: bufwriter size
const (
	_ = 1 << (10 * iota)
	KB
	MB
	GB
	TB
	PB
	EB
)

// Writer is a rotate bfw
type Writer struct {
	conf Conf
	mu   sync.Mutex
	fw   fileWriter
}

func NewWriter(conf Conf) (*Writer, error) {
	w := Writer{conf: conf}

	if conf.UseGzip == true {
		w.fw = &gzfw{bufsize: conf.BufferSize}
		conf.Suffix = conf.Suffix + ".gz"
		w.conf = conf
	} else if conf.BufferSize > 0 {
		w.fw = &bfw{bufsize: conf.BufferSize}
	} else {
		w.fw = &fw{}
	}

	if err := w.fw.Init(w.conf.newFilename()); err != nil {
		return nil, err
	}

	return &w, nil
}

func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.fw.Close(); err != nil {
		return err
	}
	return nil
}

func (w *Writer) checkRotate(p int) error {
	if w.fw.Size()+p > w.conf.MaxFilesize {
		err := w.fw.Close()
		if err != nil {
			println(1, err.Error())
		}

		err = w.fw.Init(w.conf.newFilename())
		if err != nil {
			println(2, err.Error())
		}
	}
	return nil
}

func (w *Writer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.checkRotate(len(p))
	return w.fw.Write(p)
}
