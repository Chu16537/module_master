package mgrpcclient

import (
	"context"
	"time"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mgrpc/commongrpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UnaryRPC(ctx context.Context, req *commongrpc.UnaryRPCReq) (*commongrpc.UnaryRPCRes, *errorcode.Error) {
	// 设置请求的超时时间
	cctx, cancel := context.WithTimeout(ctx, time.Duration(h.config.TimeoutSecond)*time.Second)
	defer cancel()

	res, err := h.client.UnaryRPC(cctx, req)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			return nil, errorcode.New(errorcode.Timeout, errors.Errorf("UnaryRPC eventCode:%v timeout", req.EventCode))
		}
		return nil, errorcode.Server(errors.Errorf("UnaryRPC eventCode:%v err", err))
	}

	if res.ErrorCode != errorcode.SuccessCode {
		return nil, errorcode.New(int(res.ErrorCode), errors.Errorf(res.Message))
	}

	return res, nil
}
