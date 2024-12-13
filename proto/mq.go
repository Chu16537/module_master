package proto

import "github.com/Chu16537/module_master/errorcode"

const (
	GameServer   = "game_server"
	PlayerServer = "player_server"
	// table & game
	TableToGameCreateTable = 1
	TableToGameDelTable    = 2
	TableToGameUpdateTable = 3
)

type MQSubData struct {
	SequenceID uint64 `json:"sequence_id"`
	Data       []byte `json:"data"`
}

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
		ErrorCode: errorcode.SuccessCode,
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
