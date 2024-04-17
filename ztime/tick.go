package ztime

import (
	"time"
)

const (
	Day_Sceond = 86400 // 一天秒數
)

// 每 tick 執行事件
func RunTick(interval time.Duration, f func(tick *time.Ticker)) {
	tick := time.NewTicker(interval)
	defer tick.Stop()
	f(tick)
}

// 取得時間格式 2006_01_02
func GetTimeFormatYMD(t time.Time) string {
	return t.Format("2006_01_02")
}

// 取得時間格式 2006_01_02
func GetIntFormatYMD(t int64) string {
	return time.Unix(t, 0).UTC().Format("2006_01_02")
}

// 取得 uxin
func GetUnix(t time.Time) int64 {
	return t.UTC().Unix()
}

// 生成日期范围
func GetDateRange(startTimeUnix, endTimeUnix int64) []string {
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
