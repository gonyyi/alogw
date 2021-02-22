package arotw

import "os"

// SetKeepLogs will define number of files to keep. If set to a negative number,
// it will keep everything. If 0, it will not keep previous logs.
func (w *Writer) SetKeepLogs(numOfLogsToKeep int) error {
	w.dwLog(mf_conf_setKeepLogs_dOld_dNew.Formatb(w.numKeepLogs, numOfLogsToKeep))
	w.numKeepLogs = numOfLogsToKeep
	_, err := w.delPrevLogs(w.numKeepLogs)
	return err
}

// SetMaxSize will set desired filesize per each log file.
// If same name of log to be used, it will be increased for that file.
// Default will be 10 MB.
func (w *Writer) SetMaxSize(filesize int64) {
	w.dwLog(mf_conf_setMaxSize_dOld_dNew.Formatb(w.maxSizeTarget, filesize))
	w.maxSizeTarget = filesize
	w.maxSize = filesize
}

// SetPerm defines file/directory permission
func (w *Writer) SetPerm(perm os.FileMode) {
	w.dwLog(mf_conf_setPerm_dOld_dNew.Formatb(w.perm, perm))
	w.perm = perm
}
