package logger

import (
	"io"
)

type Adapter interface {
	Writer() io.Writer
	Output(level Level, buf []byte)
}
