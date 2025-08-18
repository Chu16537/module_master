package mredis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addrs    []string      `env:"REDIS_ADDRS"`
	Password string        `env:"REDIS_PASSWORD"`
	Timeout  time.Duration `env:"REDIS_TIMEOUT_SECOND"`
	DB       int           `env:"REDIS_DB"`
}

type Handler struct {
	ctx    context.Context
	config *Config
	client redis.UniversalClient
}

func New(ctx context.Context, config *Config) (*Handler, error) {
	h := &Handler{
		ctx:    ctx,
		config: config,
		client: redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:        config.Addrs,
			Password:     config.Password,
			DB:           config.DB,
			DialTimeout:  config.Timeout,
			ReadTimeout:  config.Timeout,
			WriteTimeout: config.Timeout,
		}),
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
		if err = h.connect(); err != nil {
			return err
		}
		if _, err := h.client.Ping(ctx).Result(); err != nil {
			return fmt.Errorf("redis ping failed after reconnect: %w", err)
		}
	}

	return nil
}

// 連線
func (h *Handler) connect() error {
	ctx, cancel := context.WithTimeout(h.ctx, 10*time.Second)
	defer cancel()

	pong, err := h.client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Println("pong", pong)

	return nil
}

// 關閉
func (h *Handler) close() {
	if h.client != nil {
		_ = h.client.Close() // 忽略 Close 的錯誤
		h.client = nil
	}
}

func (h *Handler) GetRedis() redis.UniversalClient {
	return h.client
}
