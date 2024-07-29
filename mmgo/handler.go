package mmgo

import (
	"fmt"
	"time"

	"github.com/Chu16537/module_master/mgracefulshutdown"
	"github.com/Chu16537/module_master/mtime"
)

// 創建 讀寫DB
func NewMgo(level int, readConf *Config, writeConf *Config) (rc *Handler, wc *Handler, err error) {

	ctx, _ := mgracefulshutdown.GetLevelCxt(level)

	readHandler, err := New(ctx, readConf)
	if err != nil {
		return nil, nil, err
	}

	writeHandler, err := New(ctx, writeConf)
	if err != nil {
		return nil, nil, err
	}

	mgracefulshutdown.AddshutdownFunc(level, readHandler.Done)
	mgracefulshutdown.AddshutdownFunc(level, writeHandler.Done)

	rc = readHandler
	wc = writeHandler

	// 連線存活確認
	func() {
		f := func(tick *time.Ticker) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-tick.C:
					if err := readHandler.Check(); err != nil {
						fmt.Println("mongo readHandler check err", err)
					}

					if err := writeHandler.Check(); err != nil {
						fmt.Println("mongo writeHandler check err", err)
					}
				}
			}
		}

		interval := 10 * time.Second

		go mtime.RunTick(interval, f)
	}()

	return
}
