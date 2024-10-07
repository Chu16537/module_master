package mgrpcclient

import (
	"context"
	"time"

	"github.com/Chu16537/module_master/mgrpc/commongrpc"
)

func (h *Handler) UnaryRPC(ctx context.Context, req *commongrpc.UnaryRPCReq) (*commongrpc.UnaryRPCRes, error) {
	// 设置请求的超时时间
	cctx, cancel := context.WithTimeout(ctx, time.Duration(h.config.TimeoutSecond)*time.Second)
	defer cancel()

	res, err := h.client.UnaryRPC(cctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
