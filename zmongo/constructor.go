package zmongo

import (
	"context"
	"time"

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
	opts   *options.ClientOptions
	client *mongo.Client
	db     *mongo.Database
	dbName string
}

func New(ctx context.Context, conf *Config, opts *options.ClientOptions) (*Handler, error) {
	h := &Handler{
		ctx:    ctx,
		opts:   opts,
		dbName: conf.Database,
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

	if err := h.client.Ping(ctx, nil); err != nil {
		h.close()
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

	// 建立连接
	client, err := mongo.Connect(ctx, h.opts)
	if err != nil {
		return err
	}

	h.client = client
	h.db = client.Database(h.dbName)

	return nil
}

// 關閉
func (h *Handler) close() {
	if h.client != nil {
		h.client.Disconnect(h.ctx)
	}
}

// 取得 client
func (h *Handler) GetClient() *mongo.Client {
	return h.client
}

// 取得DB
func (h *Handler) GetDB() *mongo.Database {
	return h.db
}

// 取得 ctx
func (h *Handler) GetCtx() context.Context {
	return h.ctx
}
