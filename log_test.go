package goo

import (
	"testing"
)

func TestLog(t *testing.T) {
	Log.Debug("this is debug", []string{"this is debug"}, map[string]interface{}{
		"name":   "hnatao",
		"gender": 1,
	})
}
