package alogw

import (
	"os"
	"path"
	"sync"
)

// Writer is a rotate fwBuffered
type Writer struct {
	conf *Conf
	mu   sync.RWMutex // use 'em when rotating file
	fw   fw
}

func NewWriter(conf *Conf) (*Writer, error) {
	if conf == nil {
		conf = NewConf()
	}
	w := Writer{conf: conf}

	// If buffer is set, use buffered fw (fwBuffered), otherwise use basic fw (fwBasic)
	if conf.BufSize > 0 {
		w.fw = &fwBuffered{bufsize: conf.BufSize}
	} else {
		w.fw = &fwBasic{}
	}

	newFilename := w.conf.newFilename()
	if err := w.fw.Init(newFilename); err != nil {
		return nil, err
	}
	setSymlink(path.Base(newFilename), w.conf.symlink())

	return &w, nil
}

func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.fw.Close(); err != nil {
		return err
	}

	removeSymlink(w.conf.symlink())

	if err := gzipFile(w.fw.Filename()); err != nil {
		return err
	}
	if err := os.Remove(w.fw.Filename()); err != nil {
		return err
	}

	return nil
}

func (w *Writer) checkRotate(p int64) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// If size became larger but not meant to write to same file
	if (w.fw.Size()+p) > w.conf.MaxSize && w.fw.Filename() != w.conf.newFilename() {
		err := w.fw.Close()
		if err != nil {
			return err
		}

		removeSymlink(w.conf.symlink())

		if w.conf.EnableGzip {
			filename := w.fw.Filename()
			// Running as a goroutine
			go func(f string) {
				err := gzipFile(f)
				if err != nil {
					println(err.Error())
					return
				}
				// if no error, that means file has been gzip'd well.
				// therefore, delete the old file.
				if err := os.Remove(f); err != nil {
					println(err.Error())
				}
			}(filename)
		}

		newFilename := w.conf.newFilename()
		setSymlink(path.Base(newFilename), w.conf.symlink())

		err = w.fw.Init(newFilename)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Writer) Write(p []byte) (int, error) {
	w.checkRotate(int64(len(p)))
	return w.fw.Write(p)
}
