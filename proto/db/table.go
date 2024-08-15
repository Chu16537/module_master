package db

import (
	"encoding/json"
	"fmt"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	Table_Status_Unable     = 0 // 未啟用 會長不能使用
	Table_Status_Enable     = 1 // 啟用 但不能創建房間 會長 可以更改遊戲設定 會長可以更改狀態為創建
	Table_Status_WaitCreate = 2 // 等待創建 rs 會查詢 並且創建房間
	Table_Status_CreateIng  = 3 // 創建中 rs 把查詢到的 Table_Status_WaitCreate 改為 Table_Status_CreateIng 準備創建
	Table_Status_Playing    = 4 // gs 創建成功 創建失敗改為 Table_Status_WaitCreate
	Table_Status_Trans      = 5 // 房間轉換gs  gs1>rs>gs2 gs2把狀態改為 Table_Status_Playing
)

// 遊戲等級
var GameLV = map[int]map[int]struct{}{
	1: map[int]struct{}{
		1: struct{}{},
		2: struct{}{},
		3: struct{}{},
	},
	2: map[int]struct{}{
		1: struct{}{},
		2: struct{}{},
		3: struct{}{},
	},
}

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
	MaxPlayerCount    int       `json:"max_player_count" bson:"max_player_count"`         // 最大玩家數量
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

// 驗證 牌桌的 GameSetting 是否正確
func VerifyTableGameSetting(gameID int, gameConfig []byte) *errorcode.Error {
	switch gameID {
	case 1:
		return verifyTableGame1(gameConfig)
	default:
		return errorcode.New(errorcode.Game_ID_Error, errors.New(fmt.Sprintf("not gameID:%v", gameID)))

	}
}

// 遊戲1
func verifyTableGame1(data []byte) *errorcode.Error {
	gc := &GameConfig1{}
	err := json.Unmarshal(data, gc)
	if err != nil {
		return errorcode.Server(err)
	}

	// 資料驗證

	if gc.MaxPlayerCount <= 0 {
		return errorcode.New(errorcode.Game_Setting_Error, errors.New("MaxPlayerCount error"))
	}

	if len(gc.Chips) < 1 {
		return errorcode.New(errorcode.Game_Setting_Error, errors.New("Chips error"))
	}

	if len(gc.UpperBetLimitZone) < 1 {
		return errorcode.New(errorcode.Game_Setting_Error, errors.New("UpperBetLimitZone error"))
	}

	// 這3個長度要相等
	if len(gc.UpperBetLimitZone) != len(gc.Odds) || len(gc.UpperBetLimitZone) != len(gc.EffectiveBet) || len(gc.Odds) != len(gc.EffectiveBet) {
		return errorcode.New(errorcode.Game_Setting_Error, errors.New("UpperBetLimitZone Odds EffectiveBet error"))
	}

	if gc.Rtp < 0 {
		return errorcode.New(errorcode.Game_Setting_Error, errors.New("Rtp error"))
	}

	return nil
}

type TableOpt struct {
	ID            uint64
	ClubID        uint64
	GameID        int64
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

	return filter
}
