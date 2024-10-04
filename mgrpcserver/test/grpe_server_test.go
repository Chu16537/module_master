package test_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/Chu16537/module_master/mgrpcserver"
	"github.com/Chu16537/module_master/mgrpcserver/commongrpc"
	"google.golang.org/grpc"
)

func TestServer(t *testing.T) {
	config := &mgrpcserver.Config{
		Addr:          ":50051",
		TimeoutSecond: 3,
	}

	opt := []grpc.ServerOption{}

	// cgs := CommonGrpcServer()

	// h, err := mgrpcserver.New(context.Background(), config, cgs, opt...)
	h, err := mgrpcserver.New(context.Background(), config, nil, opt...)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	err = h.Run()
	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println("start", time.Now())

	doneChan := make(chan os.Signal, 1)
	signal.Notify(doneChan, syscall.SIGINT, syscall.SIGTERM)

	<-doneChan
}

// 判斷是否 超時 or client 取消
func ctxError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		fmt.Println("client is Canceled")
		return errors.New("client is Canceled")
	case context.DeadlineExceeded:
		fmt.Println("deadline is exceeded")
		return errors.New("deadline is exceeded")
	default:
		return nil
	}
}

type CommongrpcServer struct {
}

func CommonGrpcServer() *CommongrpcServer {
	return &CommongrpcServer{}
}

// UnaryRPC implements commongrpc.CommongrpcServer.
func (c *CommongrpcServer) UnaryRPC(ctx context.Context, req *commongrpc.UnaryRPCReq) (*commongrpc.UnaryRPCRes, error) {
	fmt.Println("a UnaryRPC 1")
	res := &commongrpc.UnaryRPCRes{
		EventCode: 10000,
	}
	fmt.Println("a UnaryRPC 2")

	// time.Sleep(5 * time.Second)
	return res, nil
}
