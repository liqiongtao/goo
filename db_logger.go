package goo

import (
	"fmt"
	"xorm.io/core"
)

var dbLog = NewLogger(newFileLogger("sql"))

type DBLogger struct {
	LogLevel core.LogLevel
}

func (l DBLogger) Debug(v ...interface{}) {
	dbLog.Debug(v...)
}

func (l DBLogger) Debugf(format string, v ...interface{}) {
	dbLog.Debug(fmt.Sprintf(format, v...))
}

func (l DBLogger) Error(v ...interface{}) {
	dbLog.Error(v...)
}

func (l DBLogger) Errorf(format string, v ...interface{}) {
	dbLog.Error(fmt.Sprintf(format, v...))
}

func (l DBLogger) Info(v ...interface{}) {
	dbLog.Info(v...)
}

func (l DBLogger) Infof(format string, v ...interface{}) {
	dbLog.Info(fmt.Sprintf(format, v...))
}

func (l DBLogger) Warn(v ...interface{}) {
	dbLog.Warn(v...)
}

func (l DBLogger) Warnf(format string, v ...interface{}) {
	dbLog.Warn(fmt.Sprintf(format, v...))
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
