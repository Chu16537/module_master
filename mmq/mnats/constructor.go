package mnats

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type Config struct {
	Addr            string // 多個可以使用",區隔 "nats://127.0.0.1:4222,nats://127.0.0.1:6222,nats://127.0.0.1:8222"
	User            string
	Password        string
	CreatStreamInfo []CreatStreamInfo `yaml:"creat_stream_info"`
}

type CreatStreamInfo struct {
	Name     string
	LiveTime time.Duration `yaml:"live_time"` // 生命週期
	MaxLen   int64         `yaml:"max_len"`   // 最大長度
}

type Handler struct {
	ctx    context.Context
	config *Config
	opts   []nats.Option
	nc     *nats.Conn
	js     nats.JetStreamContext

	lock   sync.RWMutex
	subMap map[string]*nats.Subscription
}

func New(ctx context.Context, config *Config) (*Handler, error) {
	h := &Handler{
		ctx:    ctx,
		config: config,
		opts:   []nats.Option{},
		subMap: make(map[string]*nats.Subscription),
	}

	// 設定 NATS 連線選項
	if config.User != "" && config.Password != "" {
		h.opts = append(h.opts, nats.UserInfo(config.User, config.Password))
	}

	if err := h.connect(); err != nil {
		return nil, err
	}

	streamInfo, err := h.js.StreamInfo(h.config.CreatStreamInfo[0].Name)
	if err != nil {
		return nil, err
	}

	fmt.Println(streamInfo.Config.Subjects)
	// for _, v := range h.config.CreatStreamInfo {
	// 	if err := h.createStream(v.Name, v.LiveTime, v.MaxLen); err != nil {
	// 		return nil, err
	// 	}
	// }

	return h, nil
}

// 關閉
func (h *Handler) Done() {
	h.lock.Lock()
	defer h.lock.Unlock()

	for _, sub := range h.subMap {
		_ = sub.Unsubscribe()
	}
	h.subMap = nil
	h.close()
}

// 檢查連線
func (h *Handler) Check() error {
	if h.nc == nil || !h.nc.IsConnected() {
		return h.connect()
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
		nc.Close()
		return err
	}

	h.nc = nc
	h.js = js
	return nil
}

// 關閉
func (h *Handler) close() {
	if h.nc != nil {
		h.nc.Close()
		h.nc = nil
	}
}
