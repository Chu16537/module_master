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

	for i := 0; i < 10; i++ {
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

	ctx := context.Background()
	stream := "order_stream"
	data := map[string]interface{}{
		"order": "{\"order_id\":1766052419435602038,\"user_id\":99991004100,\"account\":\"tokenen4100\",\"name\":\"tokenen4100\",\"agent_id\":1708411643076815707,\"agent_account\":\"USD8900\",\"master_agent_id\":1702950422629103893,\"master_agent_account\":\"idn_master\",\"game_id\":4100,\"order_credit\":2,\"effect_credit\":2,\"payout_credit\":0,\"date\":\"2024-03-13 18:25:21\",\"currency\":\"USD\",\"bet_line\":5,\"line_number\":40,\"coin_value\":0.01,\"user_credit\":100031094,\"change_credit\":-2,\"result\":{\"game_id\":4100,\"main_game\":{\"pay_credit_total\":0,\"game_result\":[[8,8,4,4],[6,6,8,7],[5,3,3,3],[1,8,8,1],[6,5,7,9]],\"pay_line\":null,\"scatter_info\":{\"id\":[9],\"position\":[[-1,-1,-1,-1],[-1,-1,-1,-1],[-1,-1,-1,-1],[-1,-1,-1,-1],[-1,-1,-1,1]],\"amount\":1,\"multiplier\":0,\"pay_credit\":0,\"pay_rate\":0},\"wild_info\":{\"id\":[0],\"position\":[[-1,-1,-1,-1],[-1,-1,-1,-1],[-1,-1,-1,-1],[-1,-1,-1,-1],[-1,-1,-1,-1]],\"amount\":0,\"multiplier\":0,\"pay_credit\":0,\"pay_rate\":0},\"scatter_extra\":[],\"extra\":{\"game_result\":[[8,8,8,4,4,8],[7,6,6,8,7,7],[5,5,3,3,3,6],[1,1,8,8,1,1],[9,6,5,7,9,8]],\"near_win\":0,\"free_spin_times\":0}},\"get_sub_game\":false,\"sub_game\":null,\"get_jackpot\":false,\"jackpot\":{\"jackpot_id\":\"\",\"jackpot_credit\":0,\"symbol_id\":null},\"get_jackpot_increment\":false,\"jackpot_increment\":null,\"grand\":0,\"major\":0,\"minor\":0,\"mini\":0,\"user_credit\":100031094,\"bet_credit\":2,\"payout_credit\":0,\"change_credit\":-2,\"effect_credit\":2,\"buy_spin\":0,\"buy_spin_multiplier\":50,\"extra\":null},\"status\":1,\"free_spin_id\":0}",
	}

	streamID, err := r.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: data,
	}).Result()

	if err != nil {
		fmt.Println("XAdd err 1", err)
		return
	}

	fmt.Println("streamID", streamID)
}
