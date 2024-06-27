package ztime

import (
	"fmt"
	"time"
)

const (
	Day_Sceond  = 86400 // 一天秒數
	Format_YMDH = "2006_01_02 15:04:05"
	Format_YMD  = "2006_01_02"
	Format_YM   = "2006_01"
)

// 每 tick 執行事件
func RunTick(interval time.Duration, f func(tick *time.Ticker)) {
	tick := time.NewTicker(interval)
	defer tick.Stop()
	f(tick)
}

// 取得時間格式 2006_01_02
func GetTimeFormat(t time.Time, format string) string {
	return t.Format(format)
}

// 取得時間格式 2006_01_02
func GetTimeFormatYMD(t time.Time) string {
	return t.Format("2006_01")
}

// 取得時間格式 2006_01_02
func GetIntFormatYMD(t int64) string {
	return time.Unix(t, 0).UTC().Format("2006_01_02")
}

// 取得 +0 時區
func GetTimeZoneZero(t time.Time) time.Time {
	return t.UTC()
}

// 取得 指定地區時間
func GetTimeZoneString(timeZone string) time.Time {
	loc, _ := time.LoadLocation(timeZone)
	tz := time.Now().In(loc)
	return tz
}

// 取得 指定地區時間
func GetTimeZoneInt(t int) time.Time {
	zoneName := "UTC"

	if t < 0 {
		zoneName = fmt.Sprintf("%v%v", zoneName, t)
	} else {
		zoneName = fmt.Sprintf("%v+%v", zoneName, t)
	}

	location := time.FixedZone(zoneName, t*3600)
	tz := time.Now().In(location)
	return tz
}

// 生成日期范围
func GetDateRange(startTimeUnix, endTimeUnix int64) []string {
	// start 時間要比較小
	if startTimeUnix > endTimeUnix {
		endTimeUnix, startTimeUnix = startTimeUnix, endTimeUnix
	}

	st := (startTimeUnix / Day_Sceond) * Day_Sceond
	et := (endTimeUnix / Day_Sceond) * Day_Sceond

	dates := make([]string, ((et-st)/Day_Sceond)+1)

	t := st
	// 循环生成日期范围
	for i := range dates {
		dates[i] = GetIntFormatYMD(t)
		t += Day_Sceond
	}

	return dates
}

// 取得今日開始秒數
func GetDayUnixStart(t int64) int64 {
	return (t / Day_Sceond) * Day_Sceond
}

// 取得今日結束秒數
func GetDayUnixEnd(t int64) int64 {
	return ((t/Day_Sceond)+1)*Day_Sceond - 1
}
