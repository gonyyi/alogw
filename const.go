package arotw

import "github.com/gonyyi/arotw/noname/anerr"

// CONST: File size
const (
	_ = 1 << (10 * iota)
	KB
	MB
	GB
	TB
	PB
	EB
)

// MESSAGES
// TODO: unexport all
const (
	// EOF marker
	mEOF    = "\n**** EOF ****\n\n"
	mEOFLen = int64(len(mEOF))

	// Unformatted Messages
	m_starting  anerr.Msg = "arotw starting"
	m_rotating anerr.Msg = "currently rotating in progress"

	// Formatted Messages
	mf_file_exist_as_dir_name_sDir anerr.Msg = "file exists with a directory name=%s"
	mf_bad_pattern_dNumPartition   anerr.Msg = "unexpected pattern: %d"
	mf_increas_max_size_dPrev_dNew anerr.Msg = "increase maxSize=%d -> %d"
	mf_remove_previous_logs        anerr.Msg = "removed previous logs=%s"
	mf_new_symlink                 anerr.Msg = "new symlink, src=%s -> %s"
	mf_close_file                  anerr.Msg = "close file=%s, size=%d"
	mf_delPrevLogs_sErr            anerr.Msg = "delPrevLogs() -> err=%s"
	mf_conf_setKeepLogs_dOld_dNew  anerr.Msg = "SetKeepLogs=%d -> %d"
	mf_conf_setMaxSize_dOld_dNew   anerr.Msg = "SetMaxSize=%d -> %d"
	mf_conf_setPerm_dOld_dNew      anerr.Msg = "SetPerm=%d -> %d"
	mf_rotate_error_sErr           anerr.Msg = "Rotate() -> err=%s"
	mf_createDirIfNotExist_sErr    anerr.Msg = "createDirIfNotExist() -> err=%s"
)
