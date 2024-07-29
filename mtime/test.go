package mtime

import (
	"fmt"
	"time"
)

func Test() {
	interval := 100 * time.Millisecond
	go RunTick(interval, aaa)
}

func aaa(tick *time.Ticker) {
	i := 0
	for {
		select {
		case <-tick.C:
			fmt.Println(i, time.Now())
			i++
			if i >= 10 {
				return
			}
		}
	}
}
