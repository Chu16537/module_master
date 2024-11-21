package mnats_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/module_master/mnats"
	"github.com/Chu16537/module_master/proto"
)

func Test_A(t *testing.T) {
	streamName := "test_stream_name"

	config := &mnats.Config{
		Addr:        "nats://127.0.0.1:4222,nats://127.0.0.1:4223,nats://127.0.0.1:4224",
		User:        "user",
		Password:    "password",
		StreamNames: []string{streamName},
	}

	ctx := context.Background()

	h, err := mnats.New(ctx, config)
	if err != nil {
		fmt.Println("New", err)
		return
	}

	subChan := make(chan proto.MQSubData, 1024)

	subName := "test_sub_name"
	consumer := "consumer2"

	m := mnats.SubMode{
		Mode: mnats.Sub_Mode_Last,
	}
	err = h.Sub(streamName, subName, consumer, m, subChan)
	if err != nil {
		fmt.Printf("Sub:%+v \n", err)
		return
	}

	defer h.UnSub(consumer)

	go func() {
		for {
			select {
			case msg, ok := <-subChan:
				fmt.Println("!ok", ok, msg.SequenceID, string(msg.Data))
				if !ok {
					return
				}

			}
		}
	}()

	time.Sleep(5 * time.Second)

	fmt.Println("start push")
	for i := 0; i < 5; i++ {
		pubData := []byte(fmt.Sprintf("unix:%v", time.Now().Unix()))
		err := h.Pub(subName, pubData)
		if err != nil {
			fmt.Println("pub err", err)
		}

		time.Sleep(1 * time.Second)
	}

	time.Sleep(2 * time.Second)
}
