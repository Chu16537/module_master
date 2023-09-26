package znats

import (
	"context"
	"fmt"
	"gomodule/znats/proto"
	"time"
)

var (
	Addr       = "127.0.0.1:4222"
	streamName = "togsa1"
	Subname    = "a1"
)

func Test() {

	ctx := context.Background()
	cfg := Config{
		Addrs: []string{Addr},
	}

	h, err := New(ctx, cfg)
	if err != nil {
		fmt.Println("err", err)
		return
	}

	h.Sub(streamName, Subname, handler)

	time.Sleep(1 * time.Second)

	msg := &proto.Req{
		Id: 1,
	}

	h.Pub(Subname, msg)

	time.Sleep(10000 * time.Second)

}

// 間聽到事件實作
func handler(msg interface{}) {
	fmt.Println("handler", msg)
}
