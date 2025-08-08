package proto

import "github.com/chu16537/module_master/errorcode"

const (
	Stream_GameServer   = "game_server"
	Stream_PlayerServer = "player_server"
	// table & game
	TableToGameCreateTable = 1 // 創建房間
	TableToGameDelTable    = 2 // 刪除房間
	TableToGameUpdateTable = 3 // 更新房間
	TableToGameGetInfo     = 4 // 取得game server 資訊
)

// table server to game server
type TableServerToGameServer struct {
	EventCode       uint64   `json:"event_code"` // 事件編號
	TableIDs        []uint64 //
	CreateTableData []CreateTableData
}

type CreateTableData struct {
	NodeID  uint64
	TableID uint64
}

// 玩家請求 推給房間
type PlayerToTable struct {
	ReqID     string `json:"req_id"`     // 請求編號
	UserID    uint64 `json:"user_id"`    // 玩家編號
	TableID   uint64 `json:"table_id"`   // 房間編號
	EventCode uint64 `json:"event_code"` // 事件編號
	Data      []byte `json:"data"`       // 資料
}

func (p *PlayerToTable) GetRes() *TableToPlayer {
	return &TableToPlayer{
		ReqID:     p.ReqID,
		UserIDs:   []uint64{p.UserID},
		TableID:   p.TableID,
		EventCode: p.EventCode,
		ErrorCode: errorcode.Code_Success,
		Data:      []byte{},
	}
}

// 房間推送資料給玩家
type TableToPlayer struct {
	ReqID     string   `json:"req_id"`     // 請求編號
	UserIDs   []uint64 `json:"user_ids"`   // 玩家編號
	TableID   uint64   `json:"table_id"`   // 房間編號
	EventCode uint64   `json:"event_code"` // 事件編號
	Data      []byte   `json:"data"`       // 資料
	ErrorCode int      `json:"error_code"` // 錯誤碼
}
