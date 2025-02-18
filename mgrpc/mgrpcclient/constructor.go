package mgrpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/Chu16537/module_master/mgrpc/commongrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
)

type Config struct {
	Addr    string
	Timeout time.Duration
}
type Handler struct {
	ctx    context.Context
	config *Config
	conn   *grpc.ClientConn
	client commongrpc.CommongrpcClient
}

func New(cxt context.Context, conf *Config) (*Handler, error) {
	h := new(Handler)

	h.ctx = cxt
	h.config = conf

	err := h.connect()
	if err != nil {
		return nil, err
	}

	h.client = commongrpc.NewCommongrpcClient(h.conn)

	return h, nil
}

// 關閉
func (h *Handler) Done() {
	if h.conn.GetState() == connectivity.Connecting {
		h.conn.Close()
	}
}

// 檢查連線
func (h *Handler) Check() error {
	status := h.conn.GetState()
	if status != connectivity.Connecting {
		h.Done()
		if err := h.connect(); err != nil {
			return fmt.Errorf("grpc connect fail %v", h.config.Addr)
		}
	}
	return nil
}

// 取得連線
func (h *Handler) GetConn() *grpc.ClientConn {
	return h.conn
}

func (h *Handler) connect() error {
	opt := keepalive.ClientParameters{
		Time:                60 * time.Second, // 定期发送 PING 帧的时间间隔
		Timeout:             60 * time.Second, // 确认 PING 帧的超时时间
		PermitWithoutStream: true,             // 允许在没有活动流时发送心跳包
	}
	conn, err := grpc.Dial(h.config.Addr, grpc.WithInsecure(), grpc.WithKeepaliveParams(opt))
	if err != nil {
		return err
	}
	h.conn = conn
	return nil
}
