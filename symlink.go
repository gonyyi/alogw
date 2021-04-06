package alogw

import "os"

func setSymlink(filename, symlink string) error {
	return os.Symlink(filename, symlink)
}

func removeSymlink(symlink string) error {
	return os.Remove(symlink)
}
