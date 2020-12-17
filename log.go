package goo

import (
	"github.com/liqiongtao/goo/logger"
)

var Log = logger.NewFileLogger("logs/", "")

func NewFileLogger(filePath, fileName string) *logger.FileLogger {
	return logger.NewFileLogger(filePath, fileName)
}
