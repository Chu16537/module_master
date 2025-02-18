package mgrpcclient

import (
	"context"
	"time"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mgrpc/commongrpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UnaryRPC(ctx context.Context, req *commongrpc.UnaryRPCReq) (*commongrpc.UnaryRPCRes, *errorcode.Error) {
	// 设置请求的超时时间
	cctx, cancel := context.WithTimeout(ctx, h.config.Timeout*time.Second)
	defer cancel()

	res, err := h.client.UnaryRPC(cctx, req)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			return nil, errorcode.New(errorcode.Code_Timeout, errors.Errorf("UnaryRPC eventCode:%v timeout", req.EventCode))
		}
		return nil, errorcode.New(errorcode.Code_Server_Error, errors.Errorf("UnaryRPC eventCode:%v err", err))
	}

	if res.ErrorCode != errorcode.Code_Success {
		return nil, errorcode.New(int(res.ErrorCode), errors.Errorf(res.Message))
	}

	return res, nil
}

func UnaryRPC(ctx context.Context, ip string, req *commongrpc.UnaryRPCReq) (*commongrpc.UnaryRPCRes, *errorcode.Error) {
	// 连接到 gRPC 服务
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		return nil, errorcode.New(errorcode.Code_Server_Error, err)
	}
	defer conn.Close()

	// 创建一个客户端
	client := commongrpc.NewCommongrpcClient(conn)

	res, err := client.UnaryRPC(ctx, req)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			return nil, errorcode.New(errorcode.Code_Timeout, errors.Errorf("UnaryRPC eventCode:%v timeout", req.EventCode))
		}
		return nil, errorcode.New(errorcode.Code_Server_Error, errors.Errorf("UnaryRPC eventCode:%v err", err))
	}

	if res.ErrorCode != errorcode.Code_Success {
		return nil, errorcode.New(int(res.ErrorCode), errors.Errorf(res.Message))
	}

	return res, nil
}
