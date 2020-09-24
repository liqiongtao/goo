package utils

import (
	"fmt"
	"os"
	"runtime"
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

func Trace() []string {
	ts := []string{}
	for i := 1; i < 8; i++ {
		_, file, line, _ := runtime.Caller(i)
		ts = append(ts, fmt.Sprintf("%s %dL", file, line))
	}
	return ts
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
