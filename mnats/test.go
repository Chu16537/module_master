package mnats

import (
	"context"
	"fmt"
	"time"

	"github.com/Chu16537/module_master/mjson"
)

var (
	Addr       = "127.0.0.1:4222"
	streamName = "togsa1"
	Subname    = "a1"
	SubChan    = make(chan []byte)
)

type Msg struct {
	ID int
}

func Test() {

	ctx := context.Background()
	cfg := &Config{
		Addrs: []string{Addr},
	}

	h, err := New(ctx, cfg)
	if err != nil {
		fmt.Println("err", err)
		return
	}

	h.Sub(Subname, handler)
	go readSubChan()

	time.Sleep(1 * time.Second)

	msg := &Msg{
		ID: 1,
	}
	data, _ := mjson.Marshal(msg)

	h.Pub(Subname, data)
	h.Pub(Subname, data)
	h.Pub(Subname, data)
	h.Pub(Subname, data)

	time.Sleep(10000 * time.Second)

}

// 間聽到事件實作
func handler(msg []byte) {
	fmt.Println("handler", msg)
	SubChan <- msg
}

func readSubChan() {
	for v := range SubChan {
		fmt.Println("ReadSubChan", time.Now(), v)
		time.Sleep(3 * time.Second)
	}
}
