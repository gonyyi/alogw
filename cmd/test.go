package main

import (
	"fmt"
	"github.com/gonyyi/alogw"
	"time"
)

func main() {
	c := alogw.NewConf()
	c.EnableGzip = true
	c.FileAggrLevel = alogw.BySecond
	c.MaxSize = alogw.KB * 1
	c.Prefix = "./data/test"
	c.Suffix = ".log"
	w, err := alogw.NewWriter(c)
	if err != nil {
		println(err.Error())
		return
	}

	for i := 0; i < 1_000_000; i++ {
		w.Write([]byte(fmt.Sprintf("time: %s, count: %d\n", time.Now().String(), i)))
	}
	w.Close()
}
