package main

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	servers := "nats://127.0.0.1:4222,nats://127.0.0.1:4223,nats://127.0.0.1:4224"
	opts := []nats.Option{nats.UserInfo("user", "password")}
	// 連接到 NATS 群集
	nc, err := nats.Connect(servers, opts...)
	if err != nil {
		fmt.Println("a1", err)
		return
	}
	defer nc.Close()

	// 啟用 JetStream
	js, err := nc.JetStream()
	if err != nil {
		fmt.Println("a2", err)
		return
	}

	// 創建 Stream（如果它尚不存在）
	streamName := "ORDERS"
	subjectName := "ORDERS.created2"

	// 取得 stream
	sInfo, err := js.StreamInfo(streamName)
	if err != nil {
		fmt.Println("a22", err.Error())
		if err.Error() == "nats: stream not found" {
			_, err = js.AddStream(&nats.StreamConfig{
				Name:      streamName,
				Subjects:  []string{subjectName},
				MaxAge:    time.Hour,            // 設置消息存活時間為1小時
				Retention: nats.WorkQueuePolicy, // 僅保留未被消費的消息
			})
			if err != nil {
				fmt.Println("a3", err)
				return
			}
		}
		return
	}

	fmt.Println("sInfo", sInfo.Config.Subjects)

	// // 設定 Consumer 並設置 AckPolicy
	// c, err := js.AddConsumer(streamName, &nats.ConsumerConfig{
	// 	Durable:   "my_durable_consumer",
	// 	AckPolicy: nats.AckAllPolicy,
	// })
	// if err != nil {
	// 	fmt.Println("a33", err)
	// 	return
	// }

	// fmt.Println("c", c.Name)

	// // 發佈訊息到 JetStream
	// t := fmt.Sprintf("%v", time.Now().Unix())
	// msg := []byte(t)
	// _, err = js.Publish(subjectName, msg)
	// if err != nil {
	// 	fmt.Println("a4", err)
	// 	return
	// }
	// fmt.Println("Published message:", string(msg))

	// 使用推送型訂閱
	sub, err := js.Subscribe(subjectName, func(msg *nats.Msg) {
		fmt.Println("Received message:", string(msg.Data))
		msg.Ack() // 確認消息
	}, nats.Durable("aa"), nats.ManualAck())
	if err != nil {
		fmt.Println("a5", err)
		return
	}
	defer sub.Unsubscribe()

	// 等待一段時間以確保訂閱接收到訊息
	time.Sleep(2 * time.Second)
	// fmt.Println("Done receiving messages")
}
