package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Addr     string
	Database string
	Username string
	Password string
}

type Handler struct {
	ctx    context.Context
	config Config
	client *mongo.Client
	db     *mongo.Database
}

func New(ctx context.Context, cfg Config) (*Handler, error) {
	h := &Handler{
		ctx:    ctx,
		config: cfg,
	}

	if !h.connect() {
		return h, fmt.Errorf("mongodb connection failed")
	}

	return h, nil
}

// 關閉
func (h *Handler) Done() {
	h.close()
}

// 檢查連線
func (h *Handler) Check() error {
	if err := h.client.Ping(h.ctx, nil); err != nil {
		h.close()
		if !h.connect() {
			return fmt.Errorf("mongodb connection lost, and reconnect failed: %v", h.config)
		}
		return fmt.Errorf("mongodb connection lost, but reconnect succeeded")
	}
	return nil
}

// 連線
func (h *Handler) connect() bool {
	fmt.Println(h.config.Addr)
	// 设置 MongoDB 连接选项
	opts := options.Client().ApplyURI(h.config.Addr)

	// 建立连接
	client, err := mongo.Connect(h.ctx, opts)
	if err != nil {
		return false
	}

	h.client = client
	h.db = client.Database(h.config.Database)

	return true
}

// 關閉
func (h *Handler) close() {
	if h.client != nil {
		h.client.Disconnect(h.ctx)
	}
}
