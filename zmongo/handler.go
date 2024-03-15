package zmongo

import (
	"fmt"
	"time"

	"github.com/Chu16537/gomodule/zgracefulshutdown"
	"github.com/Chu16537/gomodule/ztime"
)

var readDB *Handler
var writeDB *Handler

// 創建 讀寫DB
func NewReadWriteDB(level int, readConf *Config, writeConf *Config) error {

	ctx, _ := zgracefulshutdown.GetLevelCxt(level)

	rDB, err := New(ctx, readConf)
	if err != nil {
		return err
	}

	wDB, err := New(ctx, writeConf)
	if err != nil {
		return err
	}

	zgracefulshutdown.AddshutdownFunc(level, readDB.Done)
	zgracefulshutdown.AddshutdownFunc(level, writeDB.Done)

	readDB = rDB
	writeDB = wDB

	check()

	return nil
}

// 檢查是否連線正常
func check() {
	interval := 10 * time.Second

	f := func(tick *time.Ticker) {
		for {
			select {
			case <-tick.C:
				if err := readDB.Check(); err != nil {
					fmt.Println("readDB check err", err)
				}

				if err := writeDB.Check(); err != nil {
					fmt.Println("readDB check err", err)
				}
			}
		}
	}

	go ztime.RunTick(interval, f)
}

func GetDB() (r *Handler, w *Handler) {
	return readDB, writeDB
}
