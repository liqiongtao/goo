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
		Log.WithField("name", "hnatao").Debug("1111")
	}()
	go func() {
		defer wg.Done()
		Log.Debug("2222")
	}()
	wg.Wait()
}

func TestLogDebug2(t *testing.T) {
	for i := 0; i < 10; i++ {
		Log.WithField("id", 100).
			WithField("name", "hnatao").
			WithField("info", map[string]interface{}{"user": "hnatao"}).
			WithField("likes", []string{"sing", "pingpong"}).
			Debug()
		// time.Sleep(time.Second)
	}
}
