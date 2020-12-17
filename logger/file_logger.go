package logger

import (
	"time"
)

type FileLogger struct {
	adapter *FileAdapter
}

func NewFileLogger(filePath, fileName string) *FileLogger {
	adapter := NewFileAdapter(
		FilePathOption(filePath),
		FileNameOption(fileName),
		FileDateFormatOption("20060102"),
		FileMaxAgeOption(7*24*time.Hour),
		FileMaxSizeOption(1<<28),
	)
	return &FileLogger{
		adapter: adapter,
	}
}

func (fl *FileLogger) WithField(key string, value interface{}) *Logger {
	l := &Logger{Adapter: fl.adapter}
	return l.WithField(key, value)
}

func (fl *FileLogger) Debug(args ...interface{}) {
	l := &Logger{Adapter: fl.adapter}
	l.log(DEBUG, args...)
}

func (fl *FileLogger) Info(args ...interface{}) {
	l := &Logger{Adapter: fl.adapter}
	l.log(INFO, args...)
}

func (fl *FileLogger) Warn(args ...interface{}) {
	l := &Logger{Adapter: fl.adapter}
	l.log(WARN, args...)
}

func (fl *FileLogger) Error(args ...interface{}) {
	l := &Logger{Adapter: fl.adapter}
	l.log(ERROR, args...)
}

func (fl *FileLogger) Panic(args ...interface{}) {
	l := &Logger{Adapter: fl.adapter}
	l.log(PANIC, args...)
}

func (fl *FileLogger) Fatal(args ...interface{}) {
	l := &Logger{Adapter: fl.adapter}
	l.log(FATAL, args...)
}
