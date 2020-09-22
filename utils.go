package goo

import (
	"strings"
	"time"
)

var Util = gooUtil{}

type gooUtil struct {
}

func (ut gooUtil) NonceStr() string {
	return strings.ToLower(ut.Id2Code(time.Now().UnixNano()))
}

func (gooUtil) Today() string {
	return time.Now().Format("2006-01-02")
}

func (gooUtil) NextDate(d int) string {
	return time.Now().AddDate(0, 0, d).Format("2006-01-02")
}

func (gooUtil) Ts2Date(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02")
}

func (gooUtil) Ts2DateTime(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func (gooUtil) Date2Ts(date string) int64 {
	ti, _ := time.Parse("2006-01-02", date)
	return ti.Unix()
}

func (gooUtil) DateTime2Ts(dateTime string) int64 {
	ti, _ := time.Parse("2006-01-02 15:04:05", dateTime)
	return ti.Unix()
}
