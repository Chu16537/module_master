package mredis_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chu16537/module_master/database/mredis"
)

func Test_A(t *testing.T) {
	// 创建一个上下文和配置
	ctx := context.Background()
	cfg := &mredis.Config{
		// Addrs:    []string{"127.0.0.1:6379"}, // Redis 集群地址
		// Password: "aaa",                      // Redis 密码
		Addrs:    []string{"localhost:7001", "localhost:7002", "localhost:7003", "localhost:7004", "localhost:7005", "localhost:7006"}, // Redis 集群地址
		Password: "aa",                                                                                                                 // Redis 密码
		DB:       0,                                                                                                                    // Redis 数据库索引
	}

	var h *mredis.Handler
	// 初始化 Redis 集群连接处理程序
	h, err := mredis.New(ctx, cfg)
	if err != nil {
		fmt.Println("Failed to initialize Redis:", err)
		return
	}

	nodeId, err := h.GetNode(ctx, time.Now().Unix())
	if err != nil {
		fmt.Println("GetNode err", err)
		return
	}

	fmt.Println("nodeId", nodeId)
}
