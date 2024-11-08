package proto

import "github.com/Chu16537/module_master/proto/db"

const (
	RS_GET_TABLE           = 1
	RS_UPDATE_TABLE_GAME   = 2
	RS_UPDATE_TABLE_STATUS = 3
)

type RSGetTableReq struct {
	TableOpt *db.TableOpt
	FindOpt  *db.FindOpt
}

type RSGetTableRes struct {
	Tables []*db.Table `json:"tables"`
}

type RSUpdateTableGameReq struct {
	TableOpt   *db.TableOpt
	GameConfig []byte `json:"game_config"` // 每個遊戲設定不同
}

type RSUpdateTableGameRes struct {
	Table *db.Table `json:"table"`
}

type RSUpdateTableReq struct {
	TableOpt   *db.TableOpt
	Status     int   `json:"status"`
	ExpireTime int64 `json:"expire_time"`
}

type RSUpdateTableRes struct {
}
