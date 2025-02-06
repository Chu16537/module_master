package tableserverclient

import (
	"context"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mgrpc/mgrpcclient"
	"github.com/Chu16537/module_master/proto/db"
)

type IRoomServer interface {
	// 取得牌桌
	GetTable(ctx context.Context, logTracer string, tableOpt *db.TableOpt, findOpt *db.FindOpt) ([]*db.Table, *errorcode.Error)
	// 更新牌桌 遊戲設定
	UpdateTableGame(ctx context.Context, logTracer string, tableOpt *db.TableOpt, gameConfig []byte) (*db.Table, *errorcode.Error)
	// 更新牌桌狀態
	UpdateTable(ctx context.Context, logTracer string, tableOpt *db.TableOpt, status int, expireTime int64) *errorcode.Error
}

type handler struct {
	ctx    context.Context
	client *mgrpcclient.Handler
}

var h *handler

func Init(ctx context.Context, conf *mgrpcclient.Config) error {
	clientHandler, err := mgrpcclient.New(ctx, conf)
	if err != nil {
		return err
	}

	h = &handler{
		ctx:    ctx,
		client: clientHandler,
	}

	return nil
}

func (h *handler) Done() {
	h.client.Done()
}
