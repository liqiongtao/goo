package goo

import (
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
