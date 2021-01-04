package anerr

import (
	"io"
	"os"
)

func Iferr(e error, f func()) {
	if e!=nil && f!=nil {
		f()
	}
}

func Iferrw(e error, iow io.Writer) {
	if e != nil {
		iow.Write([]byte(e.Error()))
	}
}

func Iferrp(e error) {
	if e != nil {
		println(e.Error())
	}
}

func Iferrwq(e error, iow io.Writer) {
	if e != nil {
		iow.Write([]byte(e.Error()))
		os.Exit(1)
	}
}

func Iferrpq(e error) {
	if e != nil {
		println(e.Error())
		os.Exit(1)
	}
}
