package mgrpcserver

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/Chu16537/module_master/mgrpc/test"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test() {
	config := &Config{
		Addr: ":50051",
	}
	h, err := New(context.Background(), config)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	test.RegisterRoomserverServer(h.Get(), NewRoomServer())

	go func() {
		a := h.Run()
		fmt.Println("a", a)
	}()
}

type RoomServer struct {
}

func NewRoomServer() *RoomServer {
	return &RoomServer{}
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

func (r *RoomServer) Ping(stream test.Roomserver_PingServer) error {
	for {
		err := ctxError(stream.Context())

		if err != nil {
			return err
		}

		req, err := stream.Recv()

		if err == io.EOF {
			fmt.Println("no more data")
			break
		}

		if err != nil {
			return status.Errorf(codes.Unknown, "req stream recv err")
		}

		id := req.GetReqID()
		fmt.Println("Ping", id)

		res := &test.Res{
			ResID: 123,
		}

		err = stream.Send(res)

		if err != nil {
			return status.Errorf(codes.Unknown, "res stream send err")
		}
	}

	return nil
}
