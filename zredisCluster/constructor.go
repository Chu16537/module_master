package zredisCluster

import (
	"context"
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
	config *Config
	client *redis.ClusterClient
	opts   *redis.ClusterOptions
}

func New(ctx context.Context, cfg *Config) (*Handler, error) {
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

	if err := h.connect(); err != nil {
		return nil, err
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
		if err2 := h.connect(); err2 != nil {
			return err2
		}
		return err
	}
	return nil
}

// 連線
func (h *Handler) connect() error {
	h.close()

	// 建立連線
	client := redis.NewClusterClient(h.opts)
	_, err := client.Ping(h.ctx).Result()
	if err != nil {
		return err
	}

	h.client = client
	return nil
}

// 關閉
func (h *Handler) close() {
	if h.client != nil {
		h.client.Close()
	}
}

func (h *Handler) GetClient() *redis.ClusterClient {
	return h.client
}
