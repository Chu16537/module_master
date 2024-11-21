package mgracefulshutdown_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/module_master/mgracefulshutdown"
)

func Test_A(t *testing.T) {
	// 初始化
	mgracefulshutdown.Init(&mgracefulshutdown.Config{WaitTime: 15})

	// 添加 Level 1 的任務
	mgracefulshutdown.AddTask(1)
	go func() {
		defer mgracefulshutdown.DoneTask(1)
		fmt.Println("[Level 1] 開始任務")
		time.Sleep(3 * time.Second) // 模擬任務執行
		fmt.Println("[Level 1] 任務完成")
	}()

	mgracefulshutdown.AddTask(1)
	go func() {
		defer mgracefulshutdown.DoneTask(1)
		fmt.Println("[Level 11] 開始任務")
		time.Sleep(8 * time.Second) // 模擬任務執行
		fmt.Println("[Level 11] 任務完成")
	}()

	// 添加 Level 2 的任務
	mgracefulshutdown.AddTask(2)
	go func() {
		defer mgracefulshutdown.DoneTask(2)
		fmt.Println("[Level 2] 開始任務")
		time.Sleep(10 * time.Second) // 模擬任務執行
		fmt.Println("[Level 2] 任務完成")
	}()

	mgracefulshutdown.Shutdown()

	// 添加關閉函數
	mgracefulshutdown.AddshutdownFunc(1, func() {
		fmt.Println("[Level 1] 關閉函數執行")
	})
	mgracefulshutdown.AddshutdownFunc(2, func() {
		fmt.Println("[Level 2] 關閉函數執行")
	})

	// 等待平滑關閉
	fmt.Println("服務啟動，按 Ctrl+C 開始關閉流程...")
	mgracefulshutdown.WaitDone()
	fmt.Println("服務已平滑關閉")
}
