package goo

type iLogger interface {
	init()
	output(buf []byte)
}

type hookFunc func(level, message string)

type logger struct {
	logger       iLogger
	hookFuncList []hookFunc
}

func NewLogger(adapter iLogger) *logger {
	l := &logger{
		logger: adapter,
	}
	l.logger.init()
	return l
}

func (l *logger) Debug(v ...interface{}) {
	l.output(level_debug, v...)
}

func (l *logger) Error(v ...interface{}) {
	l.output(level_error, v...)
}

func (l *logger) Success(v ...interface{}) {
	l.output(level_success, v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.output(level_warn, v...)
}

func (l *logger) Info(v ...interface{}) {
	l.output(level_info, v...)
}

func (l *logger) Use(fn hookFunc) *logger {
	l.hookFuncList = append(l.hookFuncList, fn)
	return l
}

func (l *logger) output(level int, v ...interface{}) {
	li := &logInfo{
		level: level,
		data:  v,
	}
	buf := li.Json()

	AsyncFunc(l.hook(level, string(buf)))

	l.logger.output(buf)
}

func (l *logger) hook(level int, message string) func() {
	return func() {
		for _, fn := range l.hookFuncList {
			fn(logLevelMessages[level], message)
		}
	}
}
