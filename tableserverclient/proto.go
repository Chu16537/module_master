package tableserverclient

import "github.com/chu16537/module_master/proto/db"

const (
	GET_TABLE           = 1
	UPDATE_TABLE_GAME   = 2
	UPDATE_TABLE_STATUS = 3
	CREATE_TABLE        = 4
)

type GetTableReq struct {
	TableOpt *db.TableOpt `json:"table_opt"`
	FindOpt  *db.FindOpt  `json:"find_opt"`
}

type GetTableRes struct {
	Tables []*db.Table `json:"tables"`
}

type UpdateTableGameReq struct {
	TableOpt   *db.TableOpt
	GameConfig []byte `json:"game_config"` // 每個遊戲設定不同
}

type UpdateTableGameRes struct {
	Table *db.Table `json:"table"`
}

type UpdateTableReq struct {
	TableOpt   *db.TableOpt
	Status     int   `json:"status"`
	ExpireTime int64 `json:"expire_time"`
}

type UpdateTableRes struct {
}

type CreateTableReq struct {
	PlatformID uint64 `json:"platform_id"` // 俱樂部id
	ExpireTime int64  `json:"expire_time"` // 到期時間
	GameID     int    `json:"game_id"`     // 遊戲編號
}

type CreateTableRes struct {
	Success bool   `json:"Success"`
	Msg     string `json:"msg"`
}

type GameServerDoneReq struct {
	NodeID int64
}

type GameServerDoneRes struct {
	Success bool   `json:"Success"`
	Msg     string `json:"msg"`
}
