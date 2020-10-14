package goo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type logInfo struct {
	level int
	data  []interface{}
}

func (li *logInfo) Json() []byte {
	info := map[string]interface{}{
		"datetime": time.Now().Format("2006-01-02 15:04:05"),
		"level":    logLevelMessages[li.level],
		"context":  li.data,
	}

	if li.level == level_error || li.level == level_warn {
		trace := []string{}
		for i := 4; i < 8; i++ {
			_, file, line, _ := runtime.Caller(i)
			if file == "" || strings.Index(file, "runtime") > 0 {
				continue
			}
			f, _ := filepath.Rel(filepath.Dir(file), file)
			trace = append(trace, fmt.Sprintf("%s[%dL]", f, line))
		}
		info["trace"] = trace
	}

	bf := &bytes.Buffer{}
	encoder := json.NewEncoder(bf)
	encoder.SetEscapeHTML(false)
	encoder.Encode(info)

	return bf.Bytes()
}
