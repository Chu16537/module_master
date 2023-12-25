package ztime

import (
	"time"
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
