package goo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Params map[string]interface{}

func (p Params) Json() []byte {
	buf, _ := json.Marshal(p)
	return buf
}

func (p Params) Xml() []byte {
	var bf bytes.Buffer
	bf.WriteString("<xml>");
	for k, v := range p {
		vv := fmt.Sprint(v)
		if vv == "" {
			continue
		}
		str := fmt.Sprintf("<%s>%s</%s>", k, vv, k)
		bf.WriteString(str)
	}
	bf.WriteString("</xml>");
	return bf.Bytes()
}

func (p Params) QueryString() string {
	data := []string{}
	for k, v := range p {
		vv := fmt.Sprint(v)
		if vv == "" {
			continue
		}
		str := fmt.Sprintf("%s=%s", k, vv)
		data = append(data, str)
	}
	sort.Strings(data)
	return strings.Join(data, "&")
}

func (p Params) Set(name string, value interface{}) Params {
	p[name] = value
	return p
}

func (p Params) GetString(name string) string {
	if _, ok := p[name]; !ok {
		return ""
	}
	value := p[name]
	switch reflect.TypeOf(value).String() {
	case "int":
		return strconv.FormatInt(int64(value.(int)), 10)
	case "int64":
		return strconv.FormatInt(value.(int64), 10)
	case "bool":
		return strconv.FormatBool(value.(bool))
	case "float64":
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	}
	return p[name].(string)
}

func (p Params) GetBool(name string) bool {
	if _, ok := p[name]; !ok {
		return false
	}
	return p[name].(bool)
}

func (p Params) GetInt(name string) int {
	if _, ok := p[name]; !ok {
		return 0
	}
	return p[name].(int)
}

func (p Params) GetInt64(name string) int64 {
	if _, ok := p[name]; !ok {
		return 0
	}
	return p[name].(int64)
}

func (p Params) GetParams(name string) Params {
	if p[name] == nil {
		return Params{}
	}
	return p[name].(Params)
}
