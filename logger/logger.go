package logger

import (
	"encoding/json"
	"fmt"
	"github.com/liqiongtao/goo/utils"
	"os"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	mu        sync.Mutex
	v         map[string]interface{}
	hooks     []func(level Level, buf []byte)
	trimPaths []string
	Adapter   Adapter
}

func (l *Logger) log(level Level, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	defer func() { l.v = map[string]interface{}{} }()

	l.WithField("level", strings.ToLower(LevelText[level]))
	l.WithField("time", time.Now().Format("2006-01-02 15:04:05"))
	l.WithField("msg", fmt.Sprint(args...))

	if level >= ERROR {
		l.Trace()
	}

	buf, _ := json.Marshal(l.v)
	buf = append(buf, []byte("\n")...)

	for _, hook := range l.hooks {
		hook(level, buf)
	}

	if l.Adapter != nil {
		l.Adapter.Output(level, buf)
	}
}

func (l *Logger) SetTrimPaths(trimPaths ...string) {
	l.trimPaths = trimPaths
}

func (l *Logger) AddHook(fn func(level Level, buf []byte)) {
	l.hooks = append(l.hooks, fn)
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	if l.v == nil {
		l.v = map[string]interface{}{}
	}
	l.v[key] = value
	return l
}

func (l *Logger) Trace() *Logger {
	arr := utils.Trace(2)
	if len(l.trimPaths) > 0 {
		for index, item := range arr {
			for _, trimPath := range l.trimPaths {
				arr[index] = strings.Replace(item, trimPath, "", -1)
			}
		}
	}
	l.WithField("trace", arr)
	return l
}

func (l *Logger) Debug(args ...interface{}) {
	l.log(DEBUG, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log(INFO, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log(WARN, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log(ERROR, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.log(PANIC, args...)
	panic(fmt.Sprint(args...))
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log(FATAL, args...)
	os.Exit(0)
}
