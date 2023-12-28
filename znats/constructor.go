package znats

import (
	"context"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
)

type Config struct {
	Addrs []string
}

type Handler struct {
	ctx    context.Context
	config *Config
	opts   *nats.Options
	nc     *nats.Conn
	js     jetstream.JetStream
}

func New(ctx context.Context, cfg *Config) (*Handler, error) {
	h := &Handler{
		ctx:    ctx,
		config: cfg,
		opts:   &nats.Options{},
	}

	if len(h.config.Addrs) > 1 {
		h.opts.Url = strings.Join(h.config.Addrs, ",")
	} else {
		h.opts.Url = h.config.Addrs[0]
	}

	if err := h.connect(); err != nil {
		return nil, errors.Wrapf(err, "connect err :%v", err.Error())
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
		if err := h.connect(); err != nil {
			return errors.Wrapf(err, "Check connect err :%v", err.Error())
		}
	}
	return nil
}

// 連線
func (h *Handler) connect() error {
	h.close()

	nc, err := nats.Connect(h.opts.Url)

	if err != nil {
		return err
	}

	js, err := jetstream.New(nc)

	if err != nil {
		return err
	}

	h.nc = nc
	h.js = js

	return nil
}

// 關閉
func (h *Handler) close() {
	if h.nc != nil || h.nc.IsConnected() {
		h.nc.Close()
	}
}
