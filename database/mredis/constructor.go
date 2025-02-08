package mredis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addrs    []string      `env:"REDIS_ADDRS" yaml:"addrs"`
	Password string        `env:"REDIS_PASSWORD" yaml:"password"`
	Timeout  time.Duration `env:"REDIS_TIMEOUT" yaml:"timeout"`
	DB       int           `env:"REDIS_DB" yaml:"db"`
}

type Handler struct {
	ctx    context.Context
	mode   int
	config *Config
	// 群集
	clientCluster *redis.ClusterClient
	optsCluster   *redis.ClusterOptions
	// 單台
	client *redis.Client
	opts   *redis.Options
	// 實作功能
	rdsCmdable redis.Cmdable
}

func New(ctx context.Context, config *Config) (*Handler, error) {
	m := Mode_Singleton
	if len(config.Addrs) > 1 {
		m = Mode_Cluster
	}

	h := &Handler{
		ctx:    ctx,
		config: config,
		mode:   m,
		optsCluster: &redis.ClusterOptions{
			Addrs:        config.Addrs,
			Password:     config.Password,
			DialTimeout:  config.Timeout,
			ReadTimeout:  config.Timeout,
			WriteTimeout: config.Timeout,
		},
		opts: &redis.Options{
			Addr:         config.Addrs[0], // 单机模式只需要一个地址
			Password:     config.Password,
			DB:           config.DB,
			DialTimeout:  config.Timeout,
			ReadTimeout:  config.Timeout,
			WriteTimeout: config.Timeout,
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

	if _, err := h.rdsCmdable.Ping(ctx).Result(); err != nil {
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

	switch h.mode {
	case Mode_Singleton:
		h.client = redis.NewClient(h.opts)
		h.rdsCmdable = h.client
	case Mode_Cluster:
		h.clientCluster = redis.NewClusterClient(h.optsCluster)
		h.rdsCmdable = h.clientCluster
	}

	_, err := h.rdsCmdable.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

// 關閉
func (h *Handler) close() {
	if h.client != nil {
		h.client.Close()
	}
	if h.clientCluster != nil {
		h.clientCluster.Close()
	}
}

func (h *Handler) GetRedis() redis.Cmdable {
	if h.rdsCmdable != nil {
		return h.rdsCmdable
	}
	return nil
}
