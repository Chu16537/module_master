package muid_test

import (
	"fmt"
	"testing"

	"github.com/Chu16537/module_master/muid"
)

func TestMain(t *testing.T) {

	muid.New(1)

	orderId := muid.CreateOrderID()
	fmt.Println("orderId", orderId)
	fmt.Println("orderId", len(orderId))

	orderId = muid.CreateOrderID()
	fmt.Println("orderId", orderId)
	fmt.Println("orderId", len(orderId))

	orderId = muid.CreateOrderID()
	fmt.Println("orderId", orderId)
	fmt.Println("orderId", len(orderId))

	// aa := make(chan int64, 1000)

	// for i := 0; i < 1000; i++ {
	// 	go func() {
	// 		aa <- h.CreateID()
	// 	}()
	// }

	// time.Sleep(3 * time.Second)

	// ss := map[int64]bool{}

	// defer func() {
	// 	fmt.Println(len(ss))
	// }()

	// for {
	// 	select {
	// 	case val, ok := <-aa:
	// 		if !ok {
	// 			// channel 已關閉且沒有更多資料
	// 			fmt.Println("Channel closed, exiting.")
	// 			return
	// 		}
	// 		fmt.Println("Received:", val)
	// 		ss[val] = true
	// 	default:
	// 		// channel 沒有資料
	// 		fmt.Println("No more data, exiting.")
	// 		return
	// 	}
	// }

}
