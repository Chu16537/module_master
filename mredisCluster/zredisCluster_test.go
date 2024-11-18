package mrediscluster_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/module_master/mrediscluster"
)

func run(h *mrediscluster.Handler) {

	ctx := context.TODO()

	// unix := time.Now().Unix()
	// fmt.Println("unix", unix)
	// gsis, err := hmrediscluster.GetGameServerRank(h, ctx, unix, true, 20)
	// if err != nil {
	// 	fmt.Println("err", err)
	// 	return
	// }

	// fmt.Println("gsis", gsis)

	// nodeIds := make([]string, len(gsis))

	// for i, v := range gsis {
	// 	nodeIds[i] = v.Member
	// }

	// fmt.Println("nodeIds", nodeIds)

	// ipMap, err := hmrediscluster.GetGameServerIP(h, ctx, nodeIds)
	// if err != nil {
	// 	fmt.Println("err", err)
	// 	return
	// }

	// fmt.Println("ipMap", ipMap)

	// unix2 := time.Now().Unix()
	// err2 := hmrediscluster.DelGameServerIP(h, ctx, unix2)
	// if err2 != nil {
	// 	fmt.Println("err", err)
	// 	return
	// }

	// var nodeId int64 = 2
	// unix := time.Now().Unix()
	// fmt.Println("unix", unix)
	// err := hmrediscluster.UpdateGameServerRank(h, ctx, "1731913864.2", nodeId, unix, "1.2.3.4")
	// if err != nil {
	// 	fmt.Println("err", err)
	// 	return
	// }

	// key := "game_server_rank"
	// gss, err := h.GetScore(ctx, key, time.Now().Unix(), 30000, true, 30)
	// if err != nil {
	// 	fmt.Println("err", err)
	// 	return
	// }

	// for i, v := range gss {
	// 	fmt.Println(i, v.Member, v.Score)
	// }

	// err = h.AddAndUpdateZset(ctx, mrediscluster.Key_Node, float64(node), "100")
	// if err != nil {
	// 	fmt.Println("err", err)
	// 	return
	// }

	// key := "testKey"

	// err := h.Set(ctx, key, "aaa", -1)
	// if err != nil {
	// 	fmt.Println("set", err)
	// }

	// s, err := h.Get(ctx, key)
	// if err != nil {
	// 	fmt.Println("get", err)
	// } else {
	// 	fmt.Println("get s", s)
	// }

	// i, err := h.Exists(ctx, key)
	// if err != nil {
	// 	fmt.Println("Exists", err)
	// } else {
	// 	fmt.Println("Exists", i)
	// }

	// err = h.Del(ctx, key)
	// if err != nil {
	// 	fmt.Println("Del", err)
	// }

	// err = h.Del(ctx, key)
	// if err != nil {
	// 	fmt.Println("Del", err)
	// }

	// i, err = h.Exists(ctx, key)
	// if err != nil {
	// 	fmt.Println("Exists", err)
	// } else {
	// 	fmt.Println("Exists", i)
	// }

	// hashKey := "testHashKey"

	// setMap := map[string]interface{}{
	// 	"a": "aa",
	// 	"s": 11,
	// 	"d": "dd",
	// }

	// err = h.HSet(ctx, hashKey, setMap)
	// if err != nil {
	// 	fmt.Println("HSet", err)
	// }

	// keys, err := h.HKeys(ctx, hashKey)
	// if err != nil {
	// 	fmt.Println("HSet", err)
	// } else {
	// 	for i, v := range keys {
	// 		fmt.Println("keys", i, v)
	// 	}
	// }

	// err = h.HDel(ctx, hashKey, "a", "s", "d")
	// if err != nil {
	// 	fmt.Println("HDel", err)
	// }

	// listKey := "testListKey"
	// err = h.SetListFromLast(ctx, listKey, "apple", "banana", "cherry")
	// if err != nil {
	// 	panic(err)
	// }

	// err = h.SetListFromLast(ctx, listKey, "durian", "elderberry")
	// if err != nil {
	// 	panic(err)
	// }

	// err = h.SetListFromFirst(ctx, listKey, "a1", "a2")
	// if err != nil {
	// 	panic(err)
	// }

	// fruits, err := h.GetList(ctx, listKey, 1, -1)
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	for i, v := range fruits {
	// 		fmt.Println("fruits", i, v)
	// 	}
	// }

	// incrKey := "incrKey"

	// for i := 0; i < 10; i++ {
	// 	i, err := h.IncrBy(ctx, incrKey, 1)
	// 	if err != nil {
	// 		fmt.Println("err", err)
	// 		return
	// 	}
	// 	fmt.Println("i", i)
	// }

	unix := time.Now().Unix()
	fmt.Println(unix)
	err := h.UpdateNode(ctx, 0, unix)
	if err != nil {
		fmt.Println(err)
		return
	}

}
func Test_A(t *testing.T) {
	// 创建一个上下文和配置
	ctx := context.Background()
	cfg := &mrediscluster.Config{
		Addrs:    []string{"localhost:7001", "localhost:7002", "localhost:7003", "localhost:7004", "localhost:7005", "localhost:7006"}, // Redis 集群地址
		Password: "aa",                                                                                                                 // Redis 密码
	}

	var h *mrediscluster.Handler
	// 初始化 Redis 集群连接处理程序
	h, err := mrediscluster.New(ctx, cfg)
	if err != nil {
		fmt.Println("Failed to initialize Redis:", err)
		return
	}

	go check(ctx, h)

	run(h)

	// 为了演示，这里只是等待一段时间然后退出
	fmt.Println("Running your application...")
	time.Sleep(2 * time.Second)
	fmt.Println("Application finished.")
}

func check(ctx context.Context, h *mrediscluster.Handler) {
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
