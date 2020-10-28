package goo

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type fileLogger struct {
	Dir          string
	Filename     string
	Perm         uint32
	MaxSize      int64
	CurrSize     int64
	CurrNum      int
	MaxDays      int
	Writer       *os.File
	fileFullName string
	fileBaseName string
	mu           sync.RWMutex
	ctx          context.Context
}

func newFileLogger(filename string) iLogger {
	return &fileLogger{
		Dir:      "logs/",
		Filename: filename,
		Perm:     0755,
		MaxSize:  300 * 1024 * 1024,
		MaxDays:  15,
		ctx:      ctx,
	}
}

func (lf *fileLogger) init() {
	matched, _ := regexp.MatchString("/$", lf.Dir)
	if !matched {
		lf.Dir += "/"
	}
	os.MkdirAll(lf.Dir, os.FileMode(lf.Perm))

	lf.createLogFile()

	fi, _ := lf.Writer.Stat()
	if fi != nil {
		lf.CurrSize = fi.Size()
	}

	filepath.Walk(lf.Dir, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		if strings.Index(info.Name(), lf.fileBaseName) == 0 {
			lf.CurrNum++
		}
		return nil
	})

	go lf.dailyRotate()
	go lf.delOldLog()
}

func (lf *fileLogger) output(buf []byte) {
	if lf.Writer == nil {
		return
	}

	lf.mu.Lock()
	defer lf.mu.Unlock()

	if lf.CurrSize > lf.MaxSize {
		lf.fileSizeRotate()
	}

	lf.CurrSize += int64(len(string(buf)))

	if _, err := lf.Writer.Write(buf); err != nil {
		log.Println(err.Error())
	}
}

func (lf *fileLogger) createLogFile() {
	lf.fileBaseName = time.Now().Format("20060102")
	if lf.Filename != "" {
		lf.fileBaseName = fmt.Sprintf("%s_%s", lf.Filename, lf.fileBaseName)
	}
	lf.fileFullName = fmt.Sprintf("%s%s.log", lf.Dir, lf.fileBaseName)

	var err error
	lf.Writer, err = os.OpenFile(lf.fileFullName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(lf.Perm))
	if err != nil {
		log.Println(err.Error())
		return
	}

	os.Chmod(lf.fileFullName, os.FileMode(lf.Perm))
}

func (lf *fileLogger) dailyRotate() {
	defer func() {
		lf.Writer.Close()
	}()

	for {
		ti := time.Now()
		y, m, d := ti.AddDate(0, 0, 1).Date()
		nextDay := time.Date(y, m, d, 0, 0, 0, 0, ti.Location())

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(nextDay.UnixNano() - ti.UnixNano() + 100)):
			lf.mu.Lock()
			lf.createLogFile()
			lf.CurrSize = 0
			lf.CurrNum = 0
			lf.mu.Unlock()
		}
	}
}

func (lf *fileLogger) fileSizeRotate() {
	lf.Writer.Close()

	lf.CurrNum++

	newFile := fmt.Sprintf("%s%s.%02d.log", lf.Dir, lf.fileBaseName, lf.CurrNum)

	if err := os.Rename(lf.fileFullName, newFile); err != nil {
		log.Println(err.Error())
		return
	}

	lf.CurrSize = 0
	lf.createLogFile()
}

func (lf *fileLogger) delOldLog() {
	filepath.Walk(lf.Dir, func(path string, info os.FileInfo, err error) error {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "Unable to delete old log '%s', error: %v\n", path, r)
			}
		}()

		if info == nil ||
			info.IsDir() ||
			!info.ModTime().Add(24 * time.Hour * time.Duration(lf.MaxDays)).Before(time.Now()) {
			return nil
		}

		return os.Remove(path)
	})
}
