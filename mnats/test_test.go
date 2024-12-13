package mnats_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/Chu16537/module_master/mnats"
	"github.com/Chu16537/module_master/mnats/hnats"
	"github.com/Chu16537/module_master/proto"
)

func Test_A(t *testing.T) {
	config := &mnats.Config{
		Addr:        "nats://127.0.0.1:4222,nats://127.0.0.1:4223,nats://127.0.0.1:4224",
		User:        "user",
		Password:    "password",
		StreamNames: []string{"stream_game_server", "stream_game_client_server"},
	}

	ctx := context.Background()

	h, err := mnats.New(ctx, config)
	if err != nil {
		fmt.Println("New", err)
		return
	}

	// h.DelStream("stream_game_server")
	// h.DelStream("stream_game_client_server")

	var (
		subGS = make(chan proto.MQSubData, 1024)
		subGC = make(chan proto.MQSubData, 1024)
		// tableID uint64 = 1
	)

	err = hnats.SubGameServer(h, subGS)
	if err != nil {
		fmt.Println("SubGameServer", err)
		return
	}

	err = hnats.SubGameClientServer(h, subGC)
	if err != nil {
		fmt.Println("SubGameClientServer", err)
		return
	}

	fmt.Println("sub success")
	// ===
	go subPlayer(h)
	go subPlayer(h)
	time.Sleep(2 * time.Second)
	go pubPlayer(h)

	// ===
	// go pubRoom(h, tableID)
	// go pubRoom(h, tableID)
	// go subRoom(h, tableID)

	time.Sleep(10 * time.Second)
}

// 模擬 房間推給玩家
func subPlayer(h *mnats.Handler) {

	subChan := make(chan proto.MQSubData, 10)

	go func() {
		err := hnats.SubPlayer(h, subChan)
		if err != nil {
			fmt.Println("SubPlayer err", err)
			return
		}

		fmt.Println("SubPlayer success")
	}()

	go func() {
		for {
			select {
			case msg, ok := <-subChan:
				fmt.Println("subPlayer !ok", ok, msg.SequenceID, string(msg.Data))
				if !ok {
					return
				}
			}
		}
	}()

}

func pubPlayer(h *mnats.Handler) {
	data := &proto.TableToPlayer{}
	for i := 0; i < 10; i++ {
		data.ReqID = strconv.Itoa(i)

		err := hnats.PubPlayer(h, data)
		if err != nil {
			fmt.Println("PubPlayer", err)
			return
		}

		fmt.Println("PubPlayer", data.ReqID)
	}
}

func subRoom(h *mnats.Handler, tableID uint64) {
	subChan := make(chan proto.MQSubData, 10)

	go func() {
		err := hnats.SubRoom(h, tableID, subChan)
		if err != nil {
			fmt.Println("SubRoom err", err)
			return
		}

		fmt.Println("SubRoom success")
	}()

	go func() {
		for {
			select {
			case msg, ok := <-subChan:
				fmt.Println("SubRoom !ok", ok, msg.SequenceID, string(msg.Data))
				if !ok {
					return
				}
			}
		}
	}()

}

func pubRoom(h *mnats.Handler, tableID uint64) {

	data := &proto.PlayerToTable{}
	for i := 0; i < 10; i++ {
		data.ReqID = strconv.Itoa(i)

		err := hnats.PubRoom(h, tableID, data)
		if err != nil {
			fmt.Println("pubRoom", err)
			return
		}

		fmt.Println("pubRoom", data.ReqID)
	}

}
