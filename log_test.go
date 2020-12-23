package goo

import (
	"fmt"
	"github.com/liqiongtao/goo/logger"
	"sync"
	"testing"
)

func TestLogDebug(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		Log.WithField("name", "hnatao").Debug("this id debug")
	}()
	go func() {
		defer wg.Done()
		Log.Info("this is info")
	}()
	wg.Wait()
}

func TestLogHook(t *testing.T) {
	Log.AddHook(func(level logger.Level, buf []byte) {
		fmt.Println("---log.hook begin-----")
		fmt.Print(logger.LevelText[level], string(buf))
		fmt.Println("---log.hook end-----")
	})

	Log.Debug("this is debug")
	Log.WithField("name", "hnatao").Info("this is info")
	Log.Warn("this is warn")
	Log.Error("this is error")
	Log.Panic("this is panic")
	Log.Fatal("this is fatal")
}

func TestLogTrace(t *testing.T) {
	Log.Error("this is error")
	Log.Trace().Debug("this is debug")
}
