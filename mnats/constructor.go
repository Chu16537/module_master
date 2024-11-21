package mnats

import (
	"context"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type Config struct {
	Addr        string // 多個可以使用",區隔 "nats://127.0.0.1:4222,nats://127.0.0.1:6222,nats://127.0.0.1:8222"
	User        string
	Password    string
	StreamNames []string
}

type Handler struct {
	ctx    context.Context
	config *Config
	opts   []nats.Option
	nc     *nats.Conn
	js     nats.JetStreamContext

	lock   sync.Mutex
	subMap map[string]*nats.Subscription
}

func New(ctx context.Context, config *Config) (*Handler, error) {
	h := &Handler{
		ctx:    ctx,
		config: config,
		lock:   sync.Mutex{},
		subMap: make(map[string]*nats.Subscription),
	}

	// 設定opt
	opts := []nats.Option{}
	if config.User != "" && config.Password != "" {
		opts = append(opts, nats.UserInfo(config.User, config.Password))
	}

	h.opts = opts

	err := h.connect()
	if err != nil {
		return nil, errors.Wrapf(err, "connect err :%v", err.Error())
	}

	for _, stream := range h.config.StreamNames {
		err = h.createStream(stream)
		if err != nil {
			return nil, errors.Wrapf(err, "createStream err :%v", err.Error())
		}
	}

	return h, nil
}

// 關閉
func (h *Handler) Done() {
	for _, v := range h.subMap {
		v.Unsubscribe()
	}

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

	nc, err := nats.Connect(h.config.Addr, h.opts...)

	if err != nil {
		return err
	}

	js, err := nc.JetStream()
	if err != nil {
		return err
	}

	h.nc = nc
	h.js = js

	return nil
}

// 關閉
func (h *Handler) close() {
	if h.nc == nil {
		return
	}

	if h.nc.IsConnected() {
		h.nc.Close()
	}
}
