package mredis

import (
	"context"
	"fmt"
	"time"
)

// 取得 node score
func (h *Handler) GetNode(ctx context.Context, unix int64) (int64, error) {
	keys := []string{Key_Node}
	values := []interface{}{unix, NodeSecond}
	nodeId, err := h.RunLua(ctx, LuaGetNode, keys, values...)
	if err != nil {
		return 0, err
	}

	return nodeId.(int64), nil
}

// 更新node 時間
// score: nodeID /  member: unix.nodeID
func (h *Handler) UpdateNode(ctx context.Context, nodeID int64) error {
	member := fmt.Sprintf("%v.%v", time.Now().Unix(), nodeID)
	err := h.UpdateZsetMember(ctx, Key_Node, float64(nodeID), member)
	if err != nil {
		return err
	}
	return nil
}
