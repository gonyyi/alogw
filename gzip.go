package alogw

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"time"
)

func gzipFile(filename string) error {
	fi, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fi.Close()

	fo, err := os.Create(filename + time.Now().Format(".20060102-150405.gz"))
	if err != nil {
		return err
	}
	defer fo.Close()

	gw := gzip.NewWriter(fo)

	defer gw.Close()

	bfi := bufio.NewReader(fi)
	buf := make([]byte, 2048)

	for {
		n, err := bfi.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		gw.Write(buf[0:n])
	}

	return nil
}
