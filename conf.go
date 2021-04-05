package alogw

import (
	"time"
)

type Conf struct {
	MaxFilesize    int
	NumFilesToKeep int
	BufferSize     int
	UseGzip        bool
	Prefix         string
	Suffix         string
}

func NewConf() Conf {
	return Conf{
		MaxFilesize:    MB,
		NumFilesToKeep: 5,
		BufferSize:     -1,
		Prefix:         "alw",
		Suffix:         ".log",
	}
}

func (c Conf) symlink() string {
	return c.Prefix + c.Suffix
}

func (c Conf) newFilename() string {
	return c.Prefix + "-" + time.Now().Format("20060102-150405") + c.Suffix
}
