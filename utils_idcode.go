package goo

import (
	"strings"
)

type IdCode struct {
	BaseNum int64
	StepNum uint
	Words   string
}

func (ic *IdCode) Code(id int64) string {
	var str string
	wLen := len(ic.Words)
	id = (ic.BaseNum + id) << ic.StepNum
	for ; id > 0; id /= int64(wLen) {
		idx := rune(id % int64(wLen))
		str = string(ic.Words[idx]) + str
	}
	return strings.ToUpper(str)
}

func (ic *IdCode) Id(code string) int64 {
	var id int64
	wLen := len(ic.Words)
	wRunes := []rune(ic.Words)
	code = strings.ToLower(code)
	for _, w := range []rune(code) {
		for j, ww := range wRunes {
			if w == ww {
				id = id*int64(wLen) + int64(j)
				break
			}
		}
	}
	return (id >> ic.StepNum) - ic.BaseNum
}

func (gooUtil) Id2Code(id int64) string {
	return (&IdCode{BaseNum: 1001, StepNum: 21, Words: "abcdefjhigklmnpqrstuvwxyz13567890"}).Code(id)
}

func (gooUtil) Code2Id(code string) int64 {
	return (&IdCode{BaseNum: 1001, StepNum: 21, Words: "abcdefjhigklmnpqrstuvwxyz13567890"}).Id(code)
}
