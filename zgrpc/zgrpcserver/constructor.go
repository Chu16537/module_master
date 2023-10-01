package zgrpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Config struct {
	Addr string
}

type Handler struct {
	config   *Config
	listener *net.Listener
	server   *grpc.Server
}

func New(conf *Config) (*Handler, error) {
	handler := new(Handler)
	handler.config = conf

	listener, err := net.Listen("tcp", handler.config.Addr)

	if err != nil {
		return nil, fmt.Errorf("grpc server Listen err", err)
	}

	handler.listener = &listener
	handler.server = grpc.NewServer()

	return handler, nil
}

// 關閉
func (h *Handler) Done() {
	if h.server != nil {
		h.server.Stop()
	}
}

// 檢查連線
func (h *Handler) Check() error {
	return nil
}

// 啟動 server
// 必須註冊完 proto
func (h *Handler) Run() error {
	return h.server.Serve(*h.listener)
}

func (h *Handler) Get() *grpc.Server {
	return h.server
}
