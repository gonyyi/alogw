package arotw

import (
	"fmt"
	"io"
	"io/ioutil"
)

func (w *Writer) SetDebugw(dw io.Writer) {
	if dw != nil {
		w.dw = dw
	}
}
func (w *Writer) dwLog(b []byte) {
	if b != nil && len(b) > 0 {
		if b[len(b)-1] != '\n' {
			b = append(b, '\n')
		}
		go w.dw.Write(b) // to prevent lock..
	}
}

func (w *Writer) dwLogf(format string, a ...interface{}) {
	if len(a) > 0 {
		w.dwLog([]byte(fmt.Sprintf(format, a...)))
	} else {
		w.dwLog([]byte(format))
	}
}

func (w *Writer) dwDiscard() {
	w.dw = ioutil.Discard
}

func (w *Writer) dwStatus() {
	w.dwLogf("arotw status -> file.MaxSizeTarget=%d, file.MaxSize=%d, statusRotating=%t, pattern=%s, "+
		"prefix=%s, timeFmt=%s, suffx=%s, symLink=%s, "+
		"numKeepLogs=%d, curSize=%d, curFilename=%s",
		w.maxSizeTarget, w.maxSize, w.statusRotating, w.pattern,
		w.prefix, w.timeFmt, w.suffix, w.symlink,
		w.numKeepLogs, w.curSize, w.curFilename)
}
