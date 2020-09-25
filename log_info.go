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

	bf := &bytes.Buffer{}
	encoder := json.NewEncoder(bf)
	encoder.SetEscapeHTML(false)
	encoder.Encode(li.data)
	buf.WriteString(fmt.Sprintf(",\"context\":%s", bf.String()))

	if li.level == level_error || li.level == level_warn {
		ts := []string{}
		for i := 4; i < 8; i++ {
			_, file, line, _ := runtime.Caller(i)
			if file == "" || strings.Index(file, "runtime") > 0 {
				continue
			}
			f, _ := filepath.Rel(filepath.Dir(file), file)
			ts = append(ts, fmt.Sprintf("%s[%dL]", f, line))
		}
		bf, _ := json.Marshal(ts)
		buf.WriteString(fmt.Sprintf(",\"trace\":%s", string(bf)))
	}

	buf.WriteString("}")
	buf.WriteString("\n")

	return buf.Bytes()
}
