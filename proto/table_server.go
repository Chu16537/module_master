package proto

// import "github.com/Chu16537/module_master/proto/db"

// const (
// 	TS_GET_TABLE           = 1
// 	TS_UPDATE_TABLE_GAME   = 2
// 	TS_UPDATE_TABLE_STATUS = 3
// 	// TS_CREATE_TABLE        = 4
// )

// type TSGetTableReq struct {
// 	TableOpt *db.TableOpt
// 	FindOpt  *db.FindOpt
// }

// type TSGetTableRes struct {
// 	Tables []*db.Table `json:"tables"`
// }

// type TSUpdateTableGameReq struct {
// 	TableOpt   *db.TableOpt
// 	GameConfig []byte `json:"game_config"` // 每個遊戲設定不同
// }

// type TSUpdateTableGameRes struct {
// 	Table *db.Table `json:"table"`
// }

// type TSUpdateTableReq struct {
// 	TableOpt   *db.TableOpt
// 	Status     int   `json:"status"`
// 	ExpireTime int64 `json:"expire_time"`
// }

// type TSUpdateTableRes struct {
// }

// type TSCreateTableReq struct {
// 	ClubID     uint64 `json:"club_id"`     // 俱樂部id
// 	ExpireTime int64  `json:"expire_time"` // 到期時間
// 	GameID     int    `json:"game_id"`     // 遊戲編號
// }

// type TSCreateTableRes struct {
// }

// type GameServerDoneReq struct {
// 	NodeID uint64
// }

// type GameServerDoneRes struct {
// }
