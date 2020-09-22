package utils

import "time"

func Today() string {
	return time.Now().Format("2006-01-02")
}

func NextDate(d int) string {
	return time.Now().AddDate(0, 0, d).Format("2006-01-02")
}

func Ts2Date(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02")
}

func Ts2DateTime(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func Date2Ts(date string) int64 {
	ti, _ := time.Parse("2006-01-02", date)
	return ti.Unix()
}

func DateTime2Ts(dateTime string) int64 {
	ti, _ := time.Parse("2006-01-02 15:04:05", dateTime)
	return ti.Unix()
}
