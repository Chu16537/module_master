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
