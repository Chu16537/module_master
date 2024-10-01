package mgrpcserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
)

type Config struct {
	Addr          string
	TimeoutSecond int
	ReStartTime   int // 重啟時間
	ReStartCount  int // 重啟次數
}

type Handler struct {
	ctx      context.Context
	config   *Config
	listener *net.Listener
	server   *grpc.Server
}

func New(ctx context.Context, config *Config) (*Handler, error) {
	// 基本判斷
	if config.Addr == "" {
		return nil, errors.New("gin new error addr nil")
	}

	h := new(Handler)
	h.ctx = ctx
	h.config = config

	listener, err := net.Listen("tcp", h.config.Addr)

	if err != nil {
		return nil, fmt.Errorf("grpc server Listen err", err)
	}

	h.listener = &listener
	h.server = grpc.NewServer()

	return h, nil
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
	errChan := make(chan error, 1)

	go func() {
		err := h.server.Serve(*h.listener)

		if err != nil {
			errChan <- err
		}
	}()

	// 等待n秒判斷是否有錯
	select {
	case err := <-errChan:
		return err

	case <-time.After(5 * time.Second):
		// 等待5秒發現沒有錯誤
		return nil
	}
}

func (h *Handler) Get() *grpc.Server {
	return h.server
}

// 循環啟動
func (h *Handler) LoopRun(count int) {
	go func() {
		if err := h.server.Serve(*h.listener); err != nil {
			// 超過次數
			if count > h.config.ReStartCount {
				panic(err)
			}

			t := time.Duration(h.config.ReStartTime) * time.Millisecond
			time.Sleep(t)
			h.LoopRun(count + 1)
		}
	}()
}
