package goo

const (
	level_debug = iota
	level_error
	level_success
	level_warn
	level_info
)

var (
	logLevelMessages = map[int]string{
		level_debug:   "DEBUG",
		level_error:   "ERROR",
		level_success: "SUCCESS",
		level_warn:    "WARN",
		level_info:    "INFO",
	}
)

var Log = NewLogger(newFileLogger(""))
