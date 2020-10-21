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
	var buf bytes.Buffer

	buf.WriteString("{")
	buf.WriteString(fmt.Sprintf("\"datetime\":\"%s\"", time.Now().Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf(",\"level\":\"%s\"", logLevelMessages[li.level]))

	cbf := &bytes.Buffer{}
	encoder := json.NewEncoder(cbf)
	encoder.SetEscapeHTML(false)
	encoder.Encode(li.data)
	cbf.Truncate(cbf.Len() - 1)
	buf.WriteString(fmt.Sprintf(",\"context\":%s", cbf.String()))

	if li.level == level_error || li.level == level_warn {
		trace := []string{}
		for i := 3; i < 8; i++ {
			_, file, line, _ := runtime.Caller(i)
			if file == "" || strings.Index(file, "runtime") > 0 {
				continue
			}
			trace = append(trace, fmt.Sprintf("%s %dL", file, line))
		}
		tbf, _ := json.Marshal(trace)
		buf.WriteString(fmt.Sprintf(",\"trace\":%s", string(tbf)))
	}

	buf.WriteString("}")
	buf.WriteString("\n")

	return buf.Bytes()
}
