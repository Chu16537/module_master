package mtime_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/module_master/mtime"
)

func Test(t *testing.T) {
	interval := 100 * time.Millisecond
	ctx := context.Background()

	a := func() {
		fmt.Println("aaa")
	}
	go mtime.RunTick(ctx, interval, a)

	time.Sleep(10 * time.Second)
}
