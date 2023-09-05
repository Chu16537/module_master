package rediscluster

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addrs    []string
	Password string // 设置密码
	// DialTimeout  time.Duration // 设置连接超时
	// ReadTimeout  time.Duration // 设置读取超时
	// WriteTimeout time.Duration // 设置写入超时
}

type Handler struct {
	ctx    context.Context
	config Config
	client *redis.ClusterClient
	opts   *redis.ClusterOptions
}

func New(ctx context.Context, cfg Config) (*Handler, error) {
	h := &Handler{
		ctx:    ctx,
		config: cfg,
		opts: &redis.ClusterOptions{
			Addrs:        cfg.Addrs,
			Password:     cfg.Password,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}

	if !h.connect() {
		return nil, fmt.Errorf("redis connection failed")
	}

	return h, nil
}

// 關閉
func (h *Handler) Done() {
	h.close()
}

// 檢查連線
func (h *Handler) Check() error {
	if _, err := h.client.Ping(h.ctx).Result(); err != nil {
		if !h.connect() {
			return fmt.Errorf("redis connection lost, and reconnect failed: %v", h.config)
		}
		return fmt.Errorf("redis connection lost, but reconnect succeeded")
	}
	return nil
}

// 連線
func (h *Handler) connect() bool {
	h.close()

	// 建立連線
	client := redis.NewClusterClient(h.opts)
	_, err := client.Ping(h.ctx).Result()
	if err != nil {
		return false
	}

	h.client = client
	return true
}

// 關閉
func (h *Handler) close() {
	if h.client != nil {
		h.client.Close()
	}
}
