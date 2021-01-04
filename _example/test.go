package main

import (
	"github.com/gonyyi/alog"
	"github.com/gonyyi/arotw"
	"github.com/gonyyi/graceful"
	"time"
)

func main() {
	t1()
}

func t1() {
	l := alog.New(nil)

	r, err := arotw.New("./log/test{-2006-0102-1504}.log", func(r2 *arotw.Writer) {
		r2.SetDebugw(l.NewWriter(alog.Linfo, 0, "RotateWriter "))
		l.SetOutput(r2)
		r2.SetMaxSize(arotw.KB * 200)
		// Delete all and from there, keep only 5
		r2.SetKeepLogs(5)
		graceful.New(func() {
			l.Close()
		})
	})

	_ = r

	if err != nil {
		println(err.Error())
		return
	}

	for i := 0; i < 10000000; i++ {
		time.Sleep(time.Microsecond * 500)
		l.Infof("Test log %d", i)
	}
}
