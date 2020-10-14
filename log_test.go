package goo

import (
	"testing"
)

func TestLog(t *testing.T) {
	Log.Error("this is debug<p>1+1=2;2-1=1;</p>", []string{"this is debug"}, map[string]interface{}{
		"name":   "hnatao",
		"gender": 1,
	})
}
