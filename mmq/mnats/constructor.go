package mnats

import (
	"context"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type Config struct {
	Addr            string // 多個可以使用",區隔 "nats://127.0.0.1:4222,nats://127.0.0.1:6222,nats://127.0.0.1:8222"
	User            string
	Password        string
	CreatStreamInfo CreatStreamInfo `yaml:"creat_stream_info"`
}

type CreatStreamInfo struct {
	Name       string
	LiveSecond time.Duration `yaml:"live_second"` // 生命週期
	MaxLen     int64         `yaml:"max_len"`     // 最大長度
}

type Handler struct {
	ctx    context.Context
	config *Config
	opts   []nats.Option
	nc     *nats.Conn
	js     nats.JetStreamContext

	lock   sync.Mutex
	subMap map[string]*nats.Subscription
	ackMap map[string]*nats.Msg
}

func New(ctx context.Context, config *Config) (*Handler, error) {
	h := &Handler{
		ctx:    ctx,
		config: config,
		lock:   sync.Mutex{},
		subMap: make(map[string]*nats.Subscription),
		ackMap: make(map[string]*nats.Msg),
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

	err = h.createStream(h.config.CreatStreamInfo.Name, h.config.CreatStreamInfo.LiveSecond, h.config.CreatStreamInfo.MaxLen)
	if err != nil {
		return nil, errors.Wrapf(err, "createStream err :%v", err.Error())
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
