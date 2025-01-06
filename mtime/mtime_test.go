package mtime_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/module_master/mtime"
)

func Test(t *testing.T) {
	ctx := context.Background()

	a := func() {
		fmt.Println("aaa")
	}
	go mtime.RunTick(ctx, 100, a)

	time.Sleep(10 * time.Second)
}
