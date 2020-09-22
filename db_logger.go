package goo

import (
	"fmt"
	"log"
	"xorm.io/core"
)

type DBLogger struct {
	LogLevel core.LogLevel
}

func (l DBLogger) Debug(v ...interface{}) {
	log.Println(v...)
}

func (l DBLogger) Debugf(format string, v ...interface{}) {
	log.Println(fmt.Sprintf(format, v...))
}

func (l DBLogger) Error(v ...interface{}) {
	log.Println(v...)
}

func (l DBLogger) Errorf(format string, v ...interface{}) {
	log.Println(fmt.Sprintf(format, v...))
}

func (l DBLogger) Info(v ...interface{}) {
	log.Println(v...)
}

func (l DBLogger) Infof(format string, v ...interface{}) {
	log.Println(fmt.Sprintf(format, v...))
}

func (l DBLogger) Warn(v ...interface{}) {
	log.Println(v...)
}

func (l DBLogger) Warnf(format string, v ...interface{}) {
	log.Println(fmt.Sprintf(format, v...))
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
