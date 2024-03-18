package zredis

import (
	"fmt"
	"time"

	"github.com/Chu16537/gomodule/zgracefulshutdown"
	"github.com/Chu16537/gomodule/ztime"
)

// 創建 讀寫DB
func NewRedis(level int, conf *Config) (rc *Handler, err error) {

	ctx, _ := zgracefulshutdown.GetLevelCxt(level)

	rc, err = New(ctx, conf)
	if err != nil {
		return nil, err
	}

	zgracefulshutdown.AddshutdownFunc(level, rc.Done)

	func() {
		interval := 10 * time.Second

		f := func(tick *time.Ticker) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-tick.C:
					if err := rc.Check(); err != nil {
						fmt.Println("redis check err", err)
					}
				}
			}
		}

		go ztime.RunTick(interval, f)
	}()

	return
}
