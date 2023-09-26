package example

import (
	"context"
	"fmt"
	"gomodule/natscluster"
	"time"
)

var (
	Addr       = "127.0.0.1:4222"
	streamName = "togsa1"
	Subname    = "a1"
)

func Example() {

	ctx := context.Background()
	cfg := natscluster.Config{
		Addrs: []string{Addr},
	}

	h, err := natscluster.New(ctx, cfg)
	if err != nil {
		fmt.Println("err", err)
		return
	}

	h.Sub(streamName, Subname, handler)

	time.Sleep(1 * time.Second)

	msg := &Req{
		Id: 1,
	}

	h.Pub(Subname, msg)

	time.Sleep(10000 * time.Second)

}

// 間聽到事件實作
func handler(msg interface{}) {
	fmt.Println("handler", msg)
}
