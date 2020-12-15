package logger

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	PANIC
	FATAL
)

var (
	LevelText = map[Level]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		PANIC: "PANIC",
		FATAL: "FATAL",
	}
)

type FileOptionName string

const (
	FILE_NAME FileOptionName = "file_name"
)
