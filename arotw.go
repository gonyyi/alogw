package arotw

import (
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"
)

// Writer is a rotate writer
type Writer struct {
	maxSizeTarget  int64 // this is initial value given from the user, and same value will be copied to maxSizeCurr.
	maxSize        int64 // when rotating failed because same file name is given, then increase maxSizeCurr
	statusRotating bool

	// File name
	pattern string // "shorty{-2006-0102-150405}.log"
	prefix  string // "shorty"
	timeFmt string // "-2006-0102-150405"
	suffix  string // ".log"
	symlink string

	// Current file
	mu          sync.Mutex
	perm        os.FileMode // file permission to be used.
	file        *os.File    // current file's os.File
	curSize     int64       // current size of file
	curFilename string      // current filename

	// Other file option
	numKeepLogs int // number of logs to keep. If -1, it will keep forever.

	dw io.Writer // this is a debug writer
}

// New will take file pattern and permission of the file(s).
// The pattern will be "([^{]*){([^}]+)}(.*)".
// Example: "test{-2006-0102-150405}.log"
// 1. Symlinked: "test.log"
// 2. Actual file: "test-2021-0101-170001.log" (assume 5:00:01 pm 1/1/2021)
func New(filePattern string, fn ...func(*Writer)) (*Writer, error) {
	filePattern = strings.TrimPrefix(filePattern, "./")
	tmp := regexp.MustCompile("([^{]*){([^}]+)}(.*)").
		FindStringSubmatch(filePattern)

	if len(tmp) == 4 {
		r := Writer{
			maxSizeTarget: 10 * MB, // default at 10MB
			maxSize:       10 * MB,
			pattern:       filePattern,
			prefix:        tmp[1],
			timeFmt:       tmp[2],
			suffix:        tmp[3],
			symlink:       tmp[1] + tmp[3],
			curSize:       0,
			curFilename:   "",
			perm:          0644,
			numKeepLogs:   -1,
			dw:            ioutil.Discard,
		}

		for _, f := range fn {
			f(&r)
		}

		r.dwStatus()

		if err := createDirIfNotExist(r.symlink, r.perm); err != nil {
			r.dwLog(mf_createDirIfNotExist_sErr.Formatb(err.Error()))
			return nil, err
		}
		if err := r.Rotate(); err != nil {
			return nil, err
		}
		r.dwLog(m_starting.Formatb())
		return &r, nil
	}
	return nil, mf_bad_pattern_dNumPartition.Formatm(len(tmp))
}

// Do method works for a script like predefined functions.
func (w *Writer) Do(fn ...func(*Writer)) {
	for _, f := range fn {
		f(w)
	}
}

// Write will write to a file, compatible with io.Writer interface.
func (w *Writer) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if n, nTotal, err := w.write(p); err != nil {
		return n, err
	} else {
		if nTotal > w.maxSize && w.statusRotating == false { // checking rotating twice is kinda redundant but speed up
			w.statusRotating = true
			err = w.rotate() // method `rotate` does not have mutex unlike `Rotate`
			w.statusRotating = false
		}
		return n, err
	}
}

// Close will close the file.
func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.close()
}

// Rotate will rotate to a new file if not exist.
// If exist, it will append.
func (w *Writer) Rotate() (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	err = w.rotate()
	if err != nil {
		w.dwLog(mf_rotate_error_sErr.Formatb(err.Error()))
	}
	return
}

// Returns symbolic link of writer
func (w *Writer) FileSymlink() string {
	return w.symlink
}

// Returns current log file
func (w *Writer) FileCurrent() string {
	return w.curFilename
}
