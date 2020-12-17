package goo

import (
	"testing"
)

func TestLogDebug(t *testing.T) {
	Log.WithField("name", "hnatao").Debug("1111")
	Log.Debug("2222")
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
