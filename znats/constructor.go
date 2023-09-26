package znats

import (
	"context"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Config struct {
	Addrs []string
}

type Handler struct {
	ctx    context.Context
	config Config
	opts   *nats.Options
	nc     *nats.Conn
	js     jetstream.JetStream

	consumeMap map[string]jetstream.ConsumeContext // 消費者
}

func New(ctx context.Context, cfg Config) (*Handler, error) {
	h := &Handler{
		ctx:        ctx,
		config:     cfg,
		opts:       &nats.Options{},
		consumeMap: make(map[string]jetstream.ConsumeContext),
	}

	if len(h.config.Addrs) > 1 {
		h.opts.Url = strings.Join(h.config.Addrs, ",")
	} else {
		h.opts.Url = h.config.Addrs[0]
	}

	if !h.connect() {
		return nil, fmt.Errorf("nats connection failed")
	}

	return h, nil
}

// 關閉
func (h *Handler) Done() {
	h.close()
}

// 檢查連線
func (h *Handler) Check() error {
	if !h.nc.IsConnected() {
		h.connect()
		return fmt.Errorf("nats connect fail")
	}
	return nil
}

// 連線
func (h *Handler) connect() bool {
	h.close()

	nc, err := nats.Connect(h.opts.Url)

	if err != nil {
		return false
	}

	js, err := jetstream.New(nc)

	if err != nil {
		return false
	}

	h.nc = nc
	h.js = js

	return true
}

// 關閉
func (h *Handler) close() {
	if h.nc != nil && h.nc.IsConnected() {
		for _, v := range h.consumeMap {
			v.Stop()
		}
		h.nc.Close()
	}
}
