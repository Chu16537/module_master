package zredisCluster

import (
	"context"
	"fmt"
	"time"
)

func Test() {
	// 创建一个上下文和配置
	ctx := context.Background()
	cfg := &Config{
		Addrs:    []string{"localhost:7001", "localhost:7002"}, // Redis 集群地址
		Password: "",                                           // Redis 密码
	}

	var handler *Handler
	for i := 0; i < 10; i++ {
		// 初始化 Redis 集群连接处理程序
		h, err := New(ctx, cfg)
		if err != nil {
			fmt.Println("Failed to initialize Redis:", err)
			continue
		}
		handler = h
	}

	if handler == nil {
		panic("redis New fail")
	}

	defer handler.Done() // 在程序结束时关闭连接

	go check(ctx, handler)

	// 为了演示，这里只是等待一段时间然后退出
	fmt.Println("Running your application...")
	time.Sleep(10 * time.Second)
	fmt.Println("Application finished.")
}

func check(ctx context.Context, h *Handler) {
	// 設定檢查秒數
	checkInterval := 2 * time.Second
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return // 上下文已完成，退出 goroutine
		case <-ticker.C:
			if err := h.Check(); err != nil {
				fmt.Println(err)
			}
		}
	}
}
