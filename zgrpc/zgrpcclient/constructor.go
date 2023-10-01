package zgrpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Config struct {
	Addr string
}
type Handler struct {
	ctx    context.Context
	config *Config
	conn   *grpc.ClientConn
}

func New(cxt context.Context, conf *Config) (*Handler, error) {
	handler := new(Handler)

	handler.ctx = cxt
	handler.config = conf

	if err := handler.connect(); err != nil {
		return nil, fmt.Errorf("grpc connect err", err)
	}

	return handler, nil
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
	fmt.Println("h.config.Addr", h.config.Addr)
	conn, err := grpc.Dial(h.config.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	h.conn = conn
	return nil
}
