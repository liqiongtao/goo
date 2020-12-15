package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FileAdapter struct {
	currFileName string
	fileHandler  *os.File
	options      map[string]interface{}
	mu           sync.RWMutex
}

func NewFileAdapter(options ...FileOption) *FileAdapter {
	fa := new(FileAdapter)
	for _, opt := range options {
		fa.SetOption(opt.Name, opt.Value)
	}
	return fa
}

func (fa *FileAdapter) SetOption(name string, value interface{}) *FileAdapter {
	fa.mu.Lock()
	defer fa.mu.Unlock()

	if fa.options == nil {
		fa.options = map[string]interface{}{}
	}
	fa.options[name] = value
	return fa
}

func (f *FileAdapter) filePath() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	filePath := f.options[optNameFilePath].(string)
	if filePath == "" {
		filePath = "logs/"
	}
	if n := len(filePath); filePath[n-1:] != "/" {
		filePath += "/"
	}
	os.MkdirAll(filePath, 0755)
	return filePath
}

func (f *FileAdapter) fileDateFormat() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	fileDateFormat := f.options[optNameFileDateFormat].(string)
	if fileDateFormat == "" {
		fileDateFormat = "20060102"
	}
	return time.Now().Format(fileDateFormat)
}

func (f *FileAdapter) fileName() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	fileName := f.options[optNameFileName].(string)
	if fileName == "" {
		fileName = f.fileDateFormat() + ".log"
	} else {
		fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "." + f.fileDateFormat() + ".log"
	}
	return f.filePath() + fileName
}

func (f *FileAdapter) fileBackName() string {
	var n int64 = 0
	files, _ := filepath.Glob(f.fileName() + "*")
	for _, fileName := range files {
		m, _ := strconv.ParseInt(filepath.Ext(fileName)[1:], 10, 64)
		if n < m {
			n = m
		}
	}
	return fmt.Sprintf("%s.%d", f.fileName(), n+1)
}

func (f *FileAdapter) maxAge() time.Duration {
	f.mu.RLock()
	defer f.mu.RUnlock()

	maxAge := f.options[optNameFileMaxAge].(time.Duration)
	if maxAge == 0 {
		return 7 * 24 * time.Hour
	}
	return maxAge
}

func (f *FileAdapter) maxSize() int64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	maxSize := f.options[optNameFileMaxSize].(int64)
	if maxSize == 0 {
		return 1 << 28
	}
	return maxSize
}

func (f *FileAdapter) dropExpireOutFile() {
	files, _ := filepath.Glob(f.fileName() + "*")
	for _, fileName := range files {
		fi, err := os.Stat(fileName)
		if err != nil {
			continue
		}
		if fi.ModTime().After(time.Now().Add(-1 * f.maxAge())) {
			continue
		}
		go os.Remove(fileName)
	}
}

func (f *FileAdapter) close() {
	if f.fileHandler == nil {
		return
	}
	f.fileHandler.Close()
	f.fileHandler = nil
}

func (f *FileAdapter) Writer() io.Writer {
	go f.dropExpireOutFile()

	fileName := f.fileName()
	if f.currFileName != fileName {
		f.currFileName = fileName
	}

	fi, err := os.Stat(fileName)
	if !os.IsNotExist(err) && fi.Size() >= f.maxSize() {
		os.Rename(fileName, f.fileBackName())
	}

	if f.fileHandler != nil {
		f.close()
	}

	fh, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil
	}

	f.fileHandler = fh

	return f.fileHandler
}

func (f *FileAdapter) Output(level Level, buf []byte) {
	if w := f.Writer(); w != nil {
		w.Write(buf)
	}
}
