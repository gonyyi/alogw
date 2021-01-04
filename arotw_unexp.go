package arotw

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func (w *Writer) write(p []byte) (n int, tot int64, err error) {
	n, err = w.file.Write(p)
	w.curSize += int64(n)
	return n, w.curSize, err
}

func (w *Writer) close() error {
	w.dwLog(mf_close_file.Formatb(w.curFilename, w.curSize+mEOFLen))
	if w.file != nil {
		w.write([]byte(mEOF))
	}
	return w.file.Close()
}

// newFilename will return current f based on the timeFmt
func (w *Writer) newFilename() string {
	return w.prefix + time.Now().Format(w.timeFmt) + w.suffix
}

// setSymlink creates symbolic link for file src to symlink dst.
func (w *Writer) setSymlink(src, symlink string) error {
	if _, err := os.Lstat(symlink); err == nil {
		if err := os.Remove(symlink); err != nil {
			return err
		}
	}
	// use absolute path for the src file
	absFn, err := filepath.Abs(src)
	if err != nil {
		return err
	}
	if err := os.Symlink(absFn, symlink); err != nil {
		return err
	}
	w.dwLog(mf_new_symlink.Formatb(src, symlink))
	return nil
}

func (w *Writer) delPrevLogs(keep int) ([]string, error) {
	var listDeleted []string
	var err error
	list, err := w.prevLogs()

	sort.Sort(sort.Reverse(sort.StringSlice(list))) // reverse sort

	for i := keep; i < len(list); i++ {
		if list[i] != w.curFilename {
			if e := os.Remove(list[i]); e != nil {
				err = fmt.Errorf("%v: %w", err, e)
			} else {
				listDeleted = append(listDeleted, list[i])
			}
		}
	}
	return listDeleted, err
}

func (w *Writer) prevLogs() ([]string, error) {
	dir, _ := filepath.Split(w.symlink)
	if dir == "" {
		dir = "./"
	}
	var r *regexp.Regexp
	{
		r = regexp.MustCompile("[0-9]+")
		r = regexp.MustCompile(w.prefix + r.ReplaceAllLiteralString(w.timeFmt, "\\d+") + w.suffix)
	}
	// rex := w.prefix + w.timeFmt +w.suffix
	var out []string
	var err error
	filepath.Walk(dir, func(path string, info os.FileInfo, e error) error {
		// if not directory and name matches the pattern,
		if err == nil && info.IsDir() == false && r.MatchString(path) {
			out = append(out, path)
		}
		if e != nil {
			err = fmt.Errorf("%v: %w", err, e)
		}
		return nil
	})
	return out, err
}

func createDirIfNotExist(fpath string, perm os.FileMode) error {
	dir := filepath.Dir(fpath)
	fi, err := os.Stat(dir)
	if err != nil {
		// NOT EXIST
		if err2 := os.MkdirAll(dir, perm); err2 != nil {
			return fmt.Errorf("%v: %w", err, err2)
		}
	}

	// if the fpath already exists, but is a fpath, then cannot create a directory.
	if fi.IsDir() == false {
		return mf_file_exist_as_dir_name_sDir.Formatm(fi.Name())
	}
	return nil
}

func (w *Writer) rotate() error {
	// GET A NEW FILE, AND IF THE NEW FILENAME IS SAME AS CURRENT FILE,
	// increase the max size by 20% and return.
	newFilename := w.newFilename()
	if w.curFilename == newFilename {
		additionalMax := w.maxSize / 5
		w.maxSize += additionalMax
		w.dwLog(mf_increas_max_size_dPrev_dNew.Formatb(w.maxSize, w.maxSize+additionalMax))
		return nil
	}

	// IF FILE NAME IS DIFFERENT, THEN SWITCH THE LOG FILE.
	w.maxSize = w.maxSizeTarget // reset maxSize. Note: maxSize is flexible.

	// if setting has numKeepLogs, delete files.
	if w.numKeepLogs > -1 { // if numKeepLogs is -1, it's keep indefinitely
		if deletedList, err := w.delPrevLogs(w.numKeepLogs); err != nil {
			w.dwLog(mf_delPrevLogs_sErr.Formatb(err.Error()))
			return err
		} else {
			if len(deletedList) > 0 {
				w.dwLog(mf_remove_previous_logs.Formatb(strings.Join(deletedList, ",")))
			}
		}
	}

	// Writing a file
	if nf, err := os.OpenFile(newFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, w.perm); err != nil {
		return err
	} else {
		// Close current file
		if w.file != nil {
			w.close()
		}

		// Place new file; set current size=0, and create a symbolic link.
		w.file = nf
		w.curFilename = newFilename
		w.curSize = 0
		if err := w.setSymlink(w.curFilename, w.FileSymlink()); err != nil {
			return err
		}
	}

	return nil
}
