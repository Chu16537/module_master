package db

import (
	"go.mongodb.org/mongo-driver/bson"
)

const (
	Table_Status_UnEnable   = 0 // 未啟用 只有admin 可以使用
	Table_Status_Enable     = 1 // 啟用 但不能創建房間 平台 可以更改遊戲設定 ,更改狀態為等待創建
	Table_Status_WaitCreate = 2 // 等待創建 ts 會查詢 並且創建房間
	Table_Status_Playing    = 3 // 創建成功 遊戲正在進行
	Table_Status_Trans      = 4 // 房間轉換
)

type Table struct {
	TableID    uint64      `json:"table_id" bson:"table_id"`       // 牌桌id
	PlatformID uint64      `json:"platform_id" bson:"platform_id"` // 平台編號
	ExpireTime int64       `json:"expire_time" bson:"expire_time"` // 到期時間
	Status     int         `json:"status" bson:"status"`           // 狀態
	GameID     int         `json:"game_id" bson:"game_id"`         // 遊戲編號
	GameConfig *GameConfig `json:"game_config" bson:"game_config"` // 每個遊戲設定不同
}

// 遊戲1
type GameConfig struct {
	EnterBalance      uint64    `json:"enter_balance" bson:"enter_balance"`               // 入場金額
	MaxPlayerCount    uint      `json:"max_player_count" bson:"max_player_count"`         // 最大玩家數量
	Chips             []uint64  `json:"chips" bson:"chips"`                               // 籌碼
	UpperBetLimitZone []uint64  `json:"upper_bet_limit_zone" bson:"upper_bet_limit_zone"` // 上限
	Odds              []float64 `json:"odds" bson:"odds"`                                 // 賠率
	EffectiveBet      []float64 `json:"effective_bet" bson:"effective_bet"`               // 有效投注
	MaxBet            uint64    `json:"max_bet" bson:"max_bet"`                           // 單局最大總下注金額 0=無上限
	Rtp               float64   `json:"rtp" bson:"rtp"`                                   // rpt
	ExtraData         string    `json:"extra_data" bson:"extra_data"`                     // 遊戲額外資料
}

type TableOpt struct {
	TableID    uint64
	PlatformID uint64
	GameID     int
	Status     []int
	ExpireTime int64
}

func (o *TableOpt) ToMgo() bson.M {
	filter := bson.M{}

	if o == nil {
		return filter
	}

	if o.TableID > 0 {
		filter["table_id"] = o.TableID
	}

	if o.PlatformID > 0 {
		filter["platform_id"] = o.PlatformID
	}

	if o.GameID > 0 {
		filter["game_id"] = o.GameID
	}

	if len(o.Status) > 0 {
		filter["status"] = bson.M{"$in": o.Status}
	}

	if o.ExpireTime > 0 {
		filter["expire_time"] = bson.M{"$lte": o.ExpireTime}
	}

	return filter
}

type TableUpdate struct {
	Status     int
	ExpireTime int64
}

func (o *TableUpdate) ToMap() map[string]interface{} {
	update := map[string]interface{}{}

	if o.Status != -1 {
		update["status"] = o.Status
	}

	if o.ExpireTime != 0 {
		update["expire_time"] = o.ExpireTime
	}

	return update
}
