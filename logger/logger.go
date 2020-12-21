package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	mu      sync.Mutex
	v       map[string]interface{}
	hooks   []func(level Level, buf []byte)
	Adapter Adapter
}

func (l *Logger) log(level Level, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	defer func() { l.v = map[string]interface{}{} }()

	l.WithField("level", strings.ToLower(LevelText[level]))
	l.WithField("time", time.Now().Format("2006-01-02 15:04:05"))
	l.WithField("msg", fmt.Sprint(args...))

	buf, _ := json.Marshal(l.v)
	buf = append(buf, []byte("\n")...)

	for _, hook := range l.hooks {
		hook(level, buf)
	}

	if l.Adapter != nil {
		l.Adapter.Output(level, buf)
	}
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
	l.WithField("trace", string(debug.Stack()))
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
