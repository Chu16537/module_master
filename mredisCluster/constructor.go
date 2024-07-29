package mredisCluster

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addrs    []string
	Password string
}

type Handler struct {
	ctx    context.Context
	client *redis.ClusterClient
	opts   *redis.ClusterOptions
}

func New(ctx context.Context, conf *Config) (*Handler, error) {

	h := &Handler{
		ctx: ctx,
		opts: &redis.ClusterOptions{
			Addrs:        conf.Addrs,
			Password:     conf.Password,
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
	ctx, cancel := context.WithTimeout(h.ctx, 10*time.Second)
	defer cancel()

	if _, err := h.client.Ping(ctx).Result(); err != nil {
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

	ctx, cancel := context.WithTimeout(h.ctx, 10*time.Second)
	defer cancel()

	// 建立連線
	client := redis.NewClusterClient(h.opts)
	_, err := client.Ping(ctx).Result()
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
