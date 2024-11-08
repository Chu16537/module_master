package mtime

import (
	"context"
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
func RunTick(ctx context.Context, interval time.Duration, f func()) {
	tick := time.NewTicker(interval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			f()
		}
	}
}

// 取得當前增減x秒後的時間
func GetNowTimeOffset(second int) time.Time {
	return time.Now().Add(time.Duration(second) * time.Second)
}

// 取得+0時間格式
func GetTimeFormatUnix(unix int64, format string) string {
	return time.Unix(unix, 0).UTC().Format(format)
}

// 取得+0時間格式
func GetTimeFormatTime(t time.Time, format string) string {
	return t.UTC().Format(format)
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
		dates[i] = GetTimeFormatUnix(t, Format_YMD)
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
