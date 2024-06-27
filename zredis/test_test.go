package zredis_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/gomodule/zredis"
	"github.com/redis/go-redis/v9"
)

func Test_A(t *testing.T) {
	// 创建一个上下文和配置
	ctx := context.Background()
	cfg := &zredis.Config{
		Addr:     "127.0.0.1:6379", // Redis 集群地址
		Password: "",               // Redis 密码
	}

	var handler *zredis.Handler
	for i := 0; i < 10; i++ {
		// 初始化 Redis 集群连接处理程序
		h, err := zredis.New(ctx, cfg)
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

	for i := 0; i < 1; i++ {
		run(handler.GetClient())
	}

	// 为了演示，这里只是等待一段时间然后退出
	fmt.Println("Running your application...")
	time.Sleep(1 * time.Second)
	fmt.Println("Application finished.")
}

func check(ctx context.Context, h *zredis.Handler) {
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

func run(r *redis.Client) {

	ctx := context.TODO()

	// 监视指定的键
	key := "mykey"

	fields := []string{"a", "d", "s", "d"}
	m, err := r.HMGet(ctx, key, fields...).Result()
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, v := range fields {
		fmt.Println(i, v, m[i])

	}

}
