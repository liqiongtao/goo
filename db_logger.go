package goo

import (
	"fmt"
	"xorm.io/core"
)

type DBLogger struct {
	LogLevel core.LogLevel
}

func (l DBLogger) Debug(v ...interface{}) {
	Log.Debug(v...)
}

func (l DBLogger) Debugf(format string, v ...interface{}) {
	Log.Debug(fmt.Sprintf(format, v...))
}

func (l DBLogger) Error(v ...interface{}) {
	Log.Error(v...)
}

func (l DBLogger) Errorf(format string, v ...interface{}) {
	Log.Error(fmt.Sprintf(format, v...))
}

func (l DBLogger) Info(v ...interface{}) {
	Log.Info(v...)
}

func (l DBLogger) Infof(format string, v ...interface{}) {
	Log.Info(fmt.Sprintf(format, v...))
}

func (l DBLogger) Warn(v ...interface{}) {
	Log.Warn(v...)
}

func (l DBLogger) Warnf(format string, v ...interface{}) {
	Log.Warn(fmt.Sprintf(format, v...))
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
