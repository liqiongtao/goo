package utils

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	EL = "\n"
)

func FILE() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}

func LINE() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}

func DIR() string {
	_, file, _, _ := runtime.Caller(1)
	return path.Dir(file) + "/"
}

func Trace(skip int) []string {
	trace := []string{}
	if skip == 0 {
		skip = 2
	}
	for i := skip; i < 12; i++ {
		_, file, line, _ := runtime.Caller(i)
		if file == "" ||
			strings.Index(file, "runtime") > 0 ||
			strings.Index(file, "src/testing") > 0 ||
			// strings.Index(file, "pkg/mod") > 0 ||
			strings.Index(file, "vendor") > 0 {
			continue
		}
		trace = append(trace, fmt.Sprintf("%s %dL", file, line))
	}
	return trace
}

func WriteToFile(filename, content string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}
