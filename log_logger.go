package goo

type iLogger interface {
	init()
	output(level int, v ...interface{})
}

type logger struct {
	logger iLogger
}

func NewLogger(adapter iLogger) *logger {
	l := &logger{
		logger: adapter,
	}
	l.logger.init()
	return l
}

func (l *logger) Debug(v ...interface{}) {
	l.logger.output(level_debug, v...)
}

func (l *logger) Error(v ...interface{}) {
	l.logger.output(level_error, v...)
}

func (l *logger) Success(v ...interface{}) {
	l.logger.output(level_success, v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.logger.output(level_warn, v...)
}

func (l *logger) Info(v ...interface{}) {
	l.logger.output(level_info, v...)
}
