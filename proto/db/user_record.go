package db

import "go.mongodb.org/mongo-driver/bson"

/*
遊戲結束寫gr只會先寫 UserID ResultBalance Info

在把 gr轉 gr day時候 把ur剩餘資料補上 並寫入 ur day
*/
type UserRecord struct {
	UserID        uint64 `json:"user_id" bson:"user_id"`               // 玩家id
	GameRecordID  string `json:"game_record_id" bson:"game_record_id"` // 局號(唯一碼)
	CreateTime    int64  `json:"create_time" bson:"create_time"`       // 創建時間(牌局結束時間)
	ClubID        uint64 `json:"club_id" bson:"club_id"`               // 俱樂部
	TableID       uint64 `json:"table_id" bson:"table_id"`             // 房間名稱
	GameID        int    `json:"game_id" bson:"game_id"`               // 遊戲編號
	GameType      int    `json:"game_type" bson:"game_type"`           // 遊戲類型
	ResultBalance int64  `json:"result_balance" bson:"result_balance"` // 輸贏結果
	GameResult    []byte `json:"game_result" bson:"game_result"`       // 遊戲結果
	Info          []byte `json:"info" bson:"info"`                     // 遊戲資料 UserRecordInfodMultiple
}

type UserRecordInfodMultiple struct {
	BetZone           []uint64  `json:"bet_zone" bson:"bet_zone"`                       // 各區總下注
	TotalBet          uint64    `json:"total_bet" bson:"total_bet"`                     // 總下注
	WinZone           []uint64  `json:"winZ_zone" bson:"winZ_zone"`                     // 各區總贏分
	TotalWin          uint64    `json:"total_win" bson:"total_win"`                     // 總贏分
	EffectiveBet      []float64 `json:"effective_bet" bson:"effective_bet"`             // 有效投注
	TotalEffectiveBet float64   `json:"total_effective_bet" bson:"total_effective_bet"` // 總有效投注
}

type UserRecordOpt struct {
	UserID        uint64
	ClubID        uint64
	TableID       uint64
	GameID        int
	GameType      int
	StartTimeUnix int64
	EndTimeUnix   int64
}

func (o *UserRecordOpt) Filter_Mgo() bson.M {
	filter := bson.M{}

	if o == nil {
		return filter
	}

	if o.UserID > 0 {
		filter["user_id"] = o.UserID
	}

	if o.ClubID > 0 {
		filter["club_id"] = o.ClubID
	}

	if o.TableID > 0 {
		filter["table_id"] = o.TableID
	}

	if o.GameID > 0 {
		filter["game_id"] = o.GameID
	}

	if o.GameType > 0 {
		filter["game_type"] = o.GameType
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

type UserRecordTotalResult struct {
	Date  string
	Total int64
}
