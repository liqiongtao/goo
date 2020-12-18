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

func (fl *FileLogger) Debug(v ...interface{}) {
	l := &Logger{Adapter: fl.adapter}
	l.Debug(v...)
}

func (fl *FileLogger) Info(v ...interface{}) {
	l := &Logger{Adapter: fl.adapter}
	l.Info(v...)
}

func (fl *FileLogger) Warn(v ...interface{}) {
	l := &Logger{Adapter: fl.adapter}
	l.Warn(v...)
}

func (fl *FileLogger) Error(v ...interface{}) {
	l := &Logger{Adapter: fl.adapter}
	l.Error(v...)
}

func (fl *FileLogger) Panic(v ...interface{}) {
	l := &Logger{Adapter: fl.adapter}
	l.Panic(v...)
}

func (fl *FileLogger) Fatal(v ...interface{}) {
	l := &Logger{Adapter: fl.adapter}
	l.Fatal(v...)
}
