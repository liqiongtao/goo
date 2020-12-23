package logger

import (
	"time"
)

type FileLogger struct {
	adapter   *FileAdapter
	trimPaths []string
	hooks     []func(level Level, buf []byte)
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

func (fl *FileLogger) SetTrimPaths(trimPaths ...string) {
	fl.trimPaths = trimPaths
}

func (fl *FileLogger) AddHook(fn func(level Level, buf []byte)) {
	fl.hooks = append(fl.hooks, fn)
}

func (fl *FileLogger) WithField(key string, value interface{}) *Logger {
	l := &Logger{Adapter: fl.adapter, hooks: fl.hooks, trimPaths: fl.trimPaths}
	return l.WithField(key, value)
}

func (fl *FileLogger) Trace() *Logger {
	l := &Logger{Adapter: fl.adapter, hooks: fl.hooks, trimPaths: fl.trimPaths}
	return l.Trace()
}

func (fl *FileLogger) Debug(v ...interface{}) {
	l := &Logger{Adapter: fl.adapter, hooks: fl.hooks, trimPaths: fl.trimPaths}
	l.Debug(v...)
}

func (fl *FileLogger) Info(v ...interface{}) {
	l := &Logger{Adapter: fl.adapter, hooks: fl.hooks, trimPaths: fl.trimPaths}
	l.Info(v...)
}

func (fl *FileLogger) Warn(v ...interface{}) {
	l := &Logger{Adapter: fl.adapter, hooks: fl.hooks, trimPaths: fl.trimPaths}
	l.Warn(v...)
}

func (fl *FileLogger) Error(v ...interface{}) {
	l := &Logger{Adapter: fl.adapter, hooks: fl.hooks, trimPaths: fl.trimPaths}
	l.Error(v...)
}

func (fl *FileLogger) Panic(v ...interface{}) {
	l := &Logger{Adapter: fl.adapter, hooks: fl.hooks, trimPaths: fl.trimPaths}
	l.Panic(v...)
}

func (fl *FileLogger) Fatal(v ...interface{}) {
	l := &Logger{Adapter: fl.adapter, hooks: fl.hooks, trimPaths: fl.trimPaths}
	l.Fatal(v...)
}
