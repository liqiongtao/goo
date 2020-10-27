package goo

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		baseDir := ""
		trace := []string{}
		for i := 3; i < 12; i++ {
			_, file, line, _ := runtime.Caller(i)
			if index := strings.Index(file, "vendor"); baseDir == "" && index > 0 {
				baseDir = file[:index]
				continue
			}
			if file == "" ||
				strings.Index(file, "runtime") > 0 ||
				strings.Index(file, "pkg/mod") > 0 ||
				strings.Index(file, "vendor") > 0 {
				continue
			}
			trace = append(trace, fmt.Sprintf("%s %dL", file, line))
		}
		for index, s := range trace {
			trace[index] = strings.Replace(s, baseDir, "", -1)
		}
		tbf, _ := json.Marshal(trace)
		buf.WriteString(fmt.Sprintf(",\"trace\":%s", string(tbf)))
	}

	buf.WriteString("}")
	buf.WriteString("\n")

	return buf.Bytes()
}
