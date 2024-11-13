package db

import (
	"go.mongodb.org/mongo-driver/bson"
)

const (
	Table_Status_Unable     = 0 // 未啟用 會長不能使用 只有admin 可以使用
	Table_Status_Enable     = 1 // 啟用 但不能創建房間 會長 可以更改遊戲設定 會長可以更改狀態為創建
	Table_Status_WaitCreate = 2 // 等待創建 rs 會查詢 並且創建房間
	Table_Status_Playing    = 3 // gs 創建成功 創建失敗改為 Table_Status_WaitCreate
	Table_Status_Trans      = 4 // 房間轉換gs  gs1>rs>gs2 gs2把狀態改為 Table_Status_Playing
)

type Table struct {
	ID         uint64 `json:"id" bson:"id"`                   // 牌桌id
	ClubID     uint64 `json:"club_id" bson:"club_id"`         // 俱樂部編號
	LV         int    `json:"lv" bson:"lv"`                   // 牌桌等級(判斷是否可以創建遊戲)
	ExpireTime int64  `json:"expire_time" bson:"expire_time"` // 到期時間
	Status     int    `json:"status" bson:"status"`           // 狀態
	GameID     int    `json:"game_id" bson:"game_id"`         // 遊戲編號
	GameConfig []byte `json:"game_config" bson:"game_config"` // 每個遊戲設定不同
}

// 遊戲1
type GameConfig1 struct {
	EnterBalance      uint64    `json:"enter_balance" bson:"enter_balance"`               // 入場金額
	MaxPlayerCount    uint      `json:"max_player_count" bson:"max_player_count"`         // 最大玩家數量
	Chips             []uint64  `json:"chips" bson:"chips"`                               // 籌碼
	UpperBetLimitZone []uint64  `json:"upper_bet_limit_zone" bson:"upper_bet_limit_zone"` // 上限
	Odds              []float64 `json:"odds" bson:"odds"`                                 // 賠率
	EffectiveBet      []float64 `json:"effective_bet" bson:"effective_bet"`               // 有效投注
	MaxBet            uint64    `json:"max_bet" bson:"max_bet"`                           // 單局最大總下注金額 0=無上限
	Rtp               float64   `json:"rtp" bson:"rtp"`                                   // rpt
}

type GameConfig2 struct {
	EnterBalance      uint64    `json:"enter_balance" bson:"enter_balance"`               // 入場金額
	MaxPlayerCount    int       `json:"max_player_count" bson:"max_player_count"`         // 最大玩家數量
	Chips             []uint64  `json:"chips" bson:"chips"`                               // 籌碼
	UpperBetLimitZone []uint64  `json:"upper_bet_limit_zone" bson:"upper_bet_limit_zone"` // 上限
	Odds              []float64 `json:"odds" bson:"odds"`                                 // 賠率
	EffectiveBet      []float64 `json:"effective_bet" bson:"effective_bet"`               // 有效投注
	MaxBet            uint64    `json:"max_bet" bson:"max_bet"`                           // 單局最大總下注金額 0=無上限
	Rtp               float64   `json:"rtp" bson:"rtp"`                                   // rpt
}

type TableOpt struct {
	ID            uint64
	ClubID        uint64
	GameID        int64
	Status        []int
	DelExpireTime int64

	DelOpt    *TableOptDel
	CreateOpt *TableOptCreate
}

type TableOptDel struct {
	Status        []int
	DelExpireTime int64
}

type TableOptCreate struct {
	Status        []int
	DelExpireTime int64
}

func (o *TableOpt) Filter_Mgo() bson.M {
	filter := bson.M{}

	if o == nil {
		return filter
	}

	if o.ID > 0 {
		filter["id"] = o.ID
	}

	if o.ClubID > 0 {
		filter["club_id"] = o.ClubID
	}

	if o.GameID > 0 {
		filter["game_id"] = o.GameID
	}

	if len(o.Status) > 0 {
		filter["status"] = bson.M{"$in": o.Status}
	}

	if o.DelExpireTime > 0 {
		filter["expire_time"] = bson.M{"$lte": o.DelExpireTime}
	}

	if o.DelOpt != nil {
		filter["$or"] = []bson.M{
			{"status": bson.M{"$in": o.DelOpt.Status}},
			{"expire_time": bson.M{"$lte": o.DelOpt.DelExpireTime}},
		}
	}

	if o.CreateOpt != nil {
		filter["$and"] = []bson.M{
			{"status": bson.M{"$in": o.CreateOpt.Status}},
			{"expire_time": bson.M{"$gt": o.CreateOpt.DelExpireTime}},
		}
	}

	return filter
}
