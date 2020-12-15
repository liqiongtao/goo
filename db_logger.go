package goo

import (
	"fmt"
	"github.com/liqiongtao/goo/logger"
	"xorm.io/core"
)

type DBLogger struct {
	LogLevel core.LogLevel
	l        *logger.FileLogger
}

func newDBLogger(logFilePath, logFileName string) *DBLogger {
	return &DBLogger{
		l: NewFileLogger(logFilePath, logFileName),
	}
}

func (l DBLogger) Debug(v ...interface{}) {
	l.l.Debug(v...)
}

func (l DBLogger) Debugf(format string, v ...interface{}) {
	l.l.Debug(fmt.Sprintf(format, v...))
}

func (l DBLogger) Error(v ...interface{}) {
	l.l.Error(v...)
}

func (l DBLogger) Errorf(format string, v ...interface{}) {
	l.l.Error(fmt.Sprintf(format, v...))
}

func (l DBLogger) Info(v ...interface{}) {
	l.l.Info(v...)
}

func (l DBLogger) Infof(format string, v ...interface{}) {
	l.l.Info(fmt.Sprintf(format, v...))
}

func (l DBLogger) Warn(v ...interface{}) {
	l.l.Warn(v...)
}

func (l DBLogger) Warnf(format string, v ...interface{}) {
	l.l.Warn(fmt.Sprintf(format, v...))
}

func (l DBLogger) Level() core.LogLevel {
	return l.LogLevel
}

func (l DBLogger) SetLevel(ll core.LogLevel) {
	l.LogLevel = ll
}

func (l DBLogger) ShowSQL(show ...bool) {
}

func (l DBLogger) IsShowSQL() bool {
	return true
}
