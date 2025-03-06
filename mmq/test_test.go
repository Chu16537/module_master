package mmq_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/module_master/mmq"
	"github.com/Chu16537/module_master/mmq/mnats"
)

var h mmq.IMQHandler

func Test_A(t *testing.T) {
	config := &mnats.Config{
		Addr:     "nats://127.0.0.1:4222,nats://127.0.0.1:4223,nats://127.0.0.1:4224",
		User:     "user",
		Password: "password",
		CreatStreamInfo: []mnats.CreatStreamInfo{
			{
				Name:       "a",
				LiveSecond: time.Hour,
				MaxLen:     1000,
			},
		},
	}

	ctx := context.Background()

	n, err := mnats.New(ctx, config)
	if err != nil {
		fmt.Println("New", err)
		return
	}

	for i := 0; i < 10; i++ {
		err := n.AddSubjects(config.CreatStreamInfo[0].Name, fmt.Sprintf("ss-%v", i))
		if err != nil {
			fmt.Println("aaa err", err)
		}
	}

	h = n

	subChan := make(chan mmq.MQSubData, 32)
	h.Sub("ss-0", "ss-0", mmq.SubMode{Mode: mmq.Sub_Mode_Last_Ack}, subChan)

	go func() {
		for v := range subChan {
			fmt.Println("sub", v.SequenceID, string(v.Data))
		}
	}()

	subChan2 := make(chan mmq.MQSubData, 32)
	h.Sub("ss-0", "ss-1", mmq.SubMode{Mode: mmq.Sub_Mode_Last_Ack}, subChan2)

	go func() {
		for v := range subChan {
			fmt.Println("sub 1", v.SequenceID, string(v.Data))
		}
	}()

	h.Pub("ss-0", []byte("hello world 1"))
	h.Pub("ss-0", []byte("hello world 2"))
	h.Pub("ss-0", []byte("hello world 3"))
	h.Pub("ss-0", []byte("hello world 4"))
	h.Pub("ss-0", []byte("hello world 5"))
	h.Pub("ss-0", []byte("hello world 6"))

	// h.DelStream("stream_game_server")
	// h.DelStream("stream_game_client_server")

	// var (
	// 	subGS = make(chan proto.MQSubData, 1024)
	// 	subGC = make(chan proto.MQSubData, 1024)
	// 	// tableID uint64 = 1
	// )

	// err = hnats.SubGameServer(h, subGS)
	// if err != nil {
	// 	fmt.Println("SubGameServer", err)
	// 	return
	// }

	// err = hnats.SubGameClientServer(h, subGC)
	// if err != nil {
	// 	fmt.Println("SubGameClientServer", err)
	// 	return
	// }

	// fmt.Println("sub success")
	// // ===
	// go subPlayer(h)
	// go subPlayer(h)
	// time.Sleep(2 * time.Second)
	// go pubPlayer(h)

	// ===
	// go pubRoom(h, tableID)
	// go pubRoom(h, tableID)
	// go subRoom(h, tableID)

	time.Sleep(3 * time.Second)
}

// // 模擬 房間推給玩家
// func subPlayer(h *mnats.Handler) {

// 	subChan := make(chan proto.MQSubData, 10)

// 	go func() {
// 		err := hnats.SubPlayer(h, subChan)
// 		if err != nil {
// 			fmt.Println("SubPlayer err", err)
// 			return
// 		}

// 		fmt.Println("SubPlayer success")
// 	}()

// 	go func() {
// 		for {
// 			select {
// 			case msg, ok := <-subChan:
// 				fmt.Println("subPlayer !ok", ok, msg.SequenceID, string(msg.Data))
// 				if !ok {
// 					return
// 				}
// 			}
// 		}
// 	}()

// }

// func pubPlayer(h *mnats.Handler) {
// 	data := &proto.TableToPlayer{}
// 	for i := 0; i < 10; i++ {
// 		data.ReqID = strconv.Itoa(i)

// 		err := hnats.PubPlayer(h, data)
// 		if err != nil {
// 			fmt.Println("PubPlayer", err)
// 			return
// 		}

// 		fmt.Println("PubPlayer", data.ReqID)
// 	}
// }

// func subRoom(h *mnats.Handler, tableID uint64) {
// 	subChan := make(chan proto.MQSubData, 10)

// 	go func() {
// 		err := hnats.SubRoom(h, tableID, subChan)
// 		if err != nil {
// 			fmt.Println("SubRoom err", err)
// 			return
// 		}

// 		fmt.Println("SubRoom success")
// 	}()

// 	go func() {
// 		for {
// 			select {
// 			case msg, ok := <-subChan:
// 				fmt.Println("SubRoom !ok", ok, msg.SequenceID, string(msg.Data))
// 				if !ok {
// 					return
// 				}
// 			}
// 		}
// 	}()

// }

// func pubRoom(h *mnats.Handler, tableID uint64) {

// 	data := &proto.PlayerToTable{}
// 	for i := 0; i < 10; i++ {
// 		data.ReqID = strconv.Itoa(i)

// 		err := hnats.PubRoom(h, tableID, data)
// 		if err != nil {
// 			fmt.Println("pubRoom", err)
// 			return
// 		}

// 		fmt.Println("pubRoom", data.ReqID)
// 	}

// }
