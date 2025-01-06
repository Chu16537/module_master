package mgrpcserver

import (
	"context"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mgrpc/commongrpc"
)

// UnaryRPC 通用請求
func (h *Handler) UnaryRPC(ctx context.Context, req *commongrpc.UnaryRPCReq) (*commongrpc.UnaryRPCRes, error) {
	ctx, cancel := context.WithTimeout(h.ctx, h.timeout)
	defer cancel()

	// 创建一个用于接收处理结果的 channel
	resChan := make(chan *commongrpc.UnaryRPCRes)

	// 使用 goroutine 异步处理 handler
	go func() {
		var res *commongrpc.UnaryRPCRes
		defer func() {
			if r := recover(); r != nil {
				res = &commongrpc.UnaryRPCRes{
					EventCode: req.EventCode,
					ErrorCode: errorcode.Code_Handler_Not_Exist,
				}
			}

			// 将结果发送到 channel
			resChan <- res
			close(resChan)
		}()

		// 调用实际的 RPC 处理逻辑
		res, _ = h.rpcHandler.UnaryRPC(ctx, req)

	}()

	// 使用 select 来处理结果或取消
	select {
	case <-ctx.Done():
		// 回傳error 會把 錯誤訊息回傳回去
		res := &commongrpc.UnaryRPCRes{
			EventCode: req.EventCode,
			ErrorCode: errorcode.Code_Timeout,
		}
		return res, nil
	case res := <-resChan:
		return res, nil
	}
}
