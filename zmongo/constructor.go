package zmongo

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
	config *Config
	client *mongo.Client
	db     *mongo.Database
	opts   *options.ClientOptions
}

func New(ctx context.Context, cfg *Config) (*Handler, error) {
	h := &Handler{
		ctx:    ctx,
		config: cfg,
		opts:   options.Client().ApplyURI(cfg.Addr),
	}

	if !h.connect() {
		return nil, fmt.Errorf("mongodb connection failed")
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
	h.close()

	// 建立连接
	client, err := mongo.Connect(h.ctx, h.opts)
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

// 取得DB
func (h *Handler) GetDB() *mongo.Database {
	return h.db
}

func (h *Handler) GetCtx() context.Context {
	return h.ctx
}
