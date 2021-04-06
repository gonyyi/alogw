package alogw

import (
	"time"
)

// FileAggregationLevel determines if same file will be used if conditions are met
type FileAggregationLevel uint8

const (
	ByNone FileAggregationLevel = iota
	BySecond
	ByMinute
	ByHour
	ByDate
	ByMonth
	ByYear
	ByNever
)

type Conf struct {
	MaxSize int64
	// MaxFiles      int
	BufSize       int
	FileAggrLevel FileAggregationLevel
	EnableGzip    bool
	Prefix        string
	Suffix        string
}

func NewConf() *Conf {
	return &Conf{
		MaxSize: MB,
		// MaxFiles:      5,
		BufSize:       -1,
		FileAggrLevel: BySecond,
		Prefix:        "alogw",
		Suffix:        ".log",
	}
}

func (c Conf) symlink() string {
	return c.Prefix + c.Suffix
}

func (c Conf) newFilename() string {
	tmp := "20060102-150405" // defualt to be BySecond
	switch c.FileAggrLevel {
	case ByNone:
		tmp = "20060102-150405.000"
	case BySecond:
		// dont do anything as already set.
	case ByMinute:
		tmp = "20060102-1504"
	case ByHour:
		tmp = "20060102-15"
	case ByDate:
		tmp = "20060102"
	case ByMonth:
		tmp = "200601"
	case ByYear:
		tmp = "2006"
	case ByNever:
		tmp = "00"
	}
	return c.Prefix + "-" + time.Now().Format(tmp) + c.Suffix
}
