package logger

import (
	"time"
)

type FileLogger struct {
	Logger
	*FileAdapter
}

func NewFileLogger(filePath, fileName string) *FileLogger {
	adapter := NewFileAdapter(
		FilePathOption(filePath),
		FileNameOption(fileName),
		FileDateFormatOption("20060102"),
		FileMaxAgeOption(7*24*time.Hour),
		FileMaxSizeOption(1<<28),
	)
	fl := &FileLogger{
		Logger:      Logger{Adapter: adapter},
		FileAdapter: adapter,
	}
	return fl
}
