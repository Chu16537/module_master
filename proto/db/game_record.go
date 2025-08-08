package db

import (
	"fmt"

	"github.com/chu16537/module_master/errorcode"
	"github.com/chu16537/module_master/mjson"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	GameType_Multiple = 1 // 多人遊戲
)

type GameRecord struct {
	GameRecordID string `json:"game_record_id" bson:"game_record_id"` // 局號(唯一碼)
	CreateTime   int64  `json:"create_time" bson:"create_time"`       // 創建時間(牌局結束時間)
	ClubID       uint64 `json:"club_id" bson:"club_id"`               // 房間名稱
	TableID      uint64 `json:"table_id" bson:"table_id"`             // 房間名稱
	GameID       int    `json:"game_id" bson:"game_id"`               // 遊戲編號
	GameType     int    `json:"game_type" bson:"game_type"`           // 遊戲類型
	GameResult   []byte `json:"game_result" bson:"game_result"`       // 遊戲結果
	Info         []byte `json:"info" bson:"info"`                     // 遊戲資料
}

// 多人遊戲資料
type GameRecordInfoMultiple struct {
	BetZone           []uint64      `json:"bet_zone" bson:"bet_zone"`                       // 各區總下注 所有玩家總和
	TotalBet          uint64        `json:"total_bet" bson:"total_bet"`                     // 總下注 所有玩家總和
	WinZone           []uint64      `json:"win_zone" bson:"win_zone"`                       // 各區總贏分 所有玩家總和
	TotalWin          uint64        `json:"total_win" bson:"total_win"`                     // 總贏分 所有玩家總和
	EffectiveBet      []float64     `json:"effective_bet" bson:"effective_bet"`             // 有效投注
	TotalEffectiveBet float64       `json:"total_effective_bet" bson:"total_effective_bet"` // 有效投注
	OddZone           []float64     `json:"odd_zone" bson:"odd_zone"`                       // 賠率
	UserRecords       []*UserRecord `json:"user_records" bson:"user_records"`               // 下注玩家資訊
}

func (g *GameRecord) GetUserRecord() ([]*UserRecord, *errorcode.Error) {
	switch g.GameType {
	case GameType_Multiple:
		info := &GameRecordInfoMultiple{}
		err := mjson.Unmarshal(g.Info, info)
		if err != nil {
			return nil, errorcode.New(errorcode.Code_Data_Unmarshal_Error, err)
		}

		return info.UserRecords, nil

	default:
		return nil, errorcode.New(errorcode.Code_Game_Not_Type, fmt.Errorf("GameRecord:%v not type :%v", g.GameRecordID, g.GameType))
	}
}

type GameRecordOpt struct {
	GameRecordID  string
	GameRecordIDs []string
	GameID        int
	GameType      int
	TableID       uint64
	StartTimeUnix int64
	EndTimeUnix   int64
}

func (o *GameRecordOpt) Filter_Mgo() bson.M {
	filter := bson.M{}

	if o == nil {
		return filter
	}

	if len(o.GameRecordIDs) > 0 {
		filter["game_record_id"] = bson.M{"$in": o.GameRecordIDs}
	} else if o.GameRecordID != "" {
		filter["game_record_id"] = o.GameRecordID
	}

	if o.GameID > 0 {
		filter["game_id"] = o.GameID
	}

	if o.GameType > 0 {
		filter["game_type"] = o.GameID
	}

	if o.TableID > 0 {
		filter["table_id"] = o.TableID
	}

	if o.StartTimeUnix > 0 {
		if o.EndTimeUnix > 0 {
			filter["create_time"] = bson.M{"$gte": o.StartTimeUnix, "$lte": o.EndTimeUnix}
		} else {
			filter["create_time"] = bson.M{"$gte": o.StartTimeUnix}
		}
	}

	return filter
}
