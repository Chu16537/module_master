package mredis

import (
	"context"
	"fmt"
	"time"

	"github.com/chu16537/module_master/mtime"
)

// 取得 node score
func (h *Handler) GetNode(ctx context.Context) (int64, error) {
	keys := []string{Key_Node}
	values := []interface{}{time.Now().Unix(), NodeSecond}
	nodeId, err := h.RunLua(ctx, LuaGetNode, keys, values...)
	if err != nil {
		return 0, err
	}

	// 自動去更新node時間
	f := func() {
		// 更新 nodeID 時間
		err := h.updateNode(h.ctx, nodeId.(int64))
		if err != nil {
			fmt.Printf("redis UpdateNode nodeId:%v err:%s \n", nodeId, err.Error())
		}
	}

	go mtime.RunTick(h.ctx, 300, f)

	return nodeId.(int64), nil
}

// 更新node 時間
// score: nodeID /  member: unix.nodeID
func (h *Handler) updateNode(ctx context.Context, nodeID int64) error {
	member := fmt.Sprintf("%v.%v", time.Now().Unix(), nodeID)
	err := h.AddAndUpdateZset(ctx, Key_Node, float64(nodeID), member)
	if err != nil {
		return err
	}
	return nil
}
