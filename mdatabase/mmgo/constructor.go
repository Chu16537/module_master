package mmgo

import (
	"context"
	"fmt"
	"time"

	"github.com/chu16537/module_master/proto/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Master *Option `yaml:"master"`
	Second *Option `yaml:"second"`
}

type Option struct {
	Addr     string `yaml:"addr"`
	Database string `yaml:"data_base"`
	Username string `yaml:"user_name"`
	Password string `yaml:"pass_word"`
}

type Handler struct {
	ctx    context.Context
	config *Config
	read   *dbOpt
	write  *dbOpt
}

type dbOpt struct {
	opts   *options.ClientOptions
	client *mongo.Client
	db     *mongo.Database
}

func New(ctx context.Context, config *Config) (*Handler, error) {
	// 基本檢查
	if config.Master == nil || config.Master.Addr == "" || config.Master.Database == "" {
		return nil, fmt.Errorf("invalid config")
	}

	//  假如read是空的 讀寫就在同一台
	if config.Second.Addr == "" {
		config.Second = config.Master
	}

	wOpt := options.Client().ApplyURI(config.Master.Addr)
	if config.Master.Username != "" {
		cred := options.Credential{
			Username: config.Master.Username,
			Password: config.Master.Password,
		}
		wOpt.SetAuth(cred)
	}

	rOpt := options.Client().ApplyURI(config.Second.Addr)
	if config.Second.Username != "" {
		cred := options.Credential{
			Username: config.Second.Username,
			Password: config.Second.Password,
		}
		rOpt.SetAuth(cred)
	}

	h := &Handler{
		ctx:    ctx,
		config: config,
		write: &dbOpt{
			opts: wOpt,
		},
		read: &dbOpt{
			opts: rOpt,
		},
	}

	if err := h.connect(); err != nil {
		return nil, err
	}

	// 初始化 count 資料表
	h.CreateCollection(ctx, db.ColName_Counters, [][]string{[]string{"col_name"}})

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

	err := h.write.client.Ping(ctx, nil)
	if err == nil {
		err = h.read.client.Ping(ctx, nil)
	}

	if err != nil {
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
	wClient, err := mongo.Connect(ctx, h.write.opts)
	if err != nil {
		return err
	}
	h.write.client = wClient
	h.write.db = h.write.client.Database(h.config.Master.Database)

	// 建立连接
	rClient, err := mongo.Connect(ctx, h.read.opts)
	if err != nil {
		return err
	}

	h.read.client = rClient
	h.read.db = h.read.client.Database(h.config.Second.Database)

	return nil
}

// 關閉
func (h *Handler) close() {
	if h.write.client != nil {
		h.write.client.Disconnect(h.ctx)
	}

	if h.read.client != nil {
		h.read.client.Disconnect(h.ctx)
	}
}

// 取得 ctx
func (h *Handler) GetCtx() context.Context {
	return h.ctx
}

// 取得 讀:讀連線 寫:寫連線
func (h *Handler) GetWR() *Handler {
	return &Handler{
		ctx:   h.ctx,
		read:  h.read,
		write: h.write,
	}
}

// 取得 讀:讀連線 寫:讀連線
func (h *Handler) GetWW() *Handler {
	return &Handler{
		ctx:   h.ctx,
		read:  h.write,
		write: h.write,
	}
}
