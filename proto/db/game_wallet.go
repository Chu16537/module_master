package db

import (
	"go.mongodb.org/mongo-driver/bson"
)

// 遊戲錢包
type GameWallet struct {
	UserID            uint64  `json:"user_id" bson:"user_id"` // 玩家名稱
	ClubID            uint64  `json:"club_id" bson:"club_id"`
	TableID           uint64  `json:"table_id" bson:"table_id"`                       // 房間名稱
	DelTime           int64   `json:"del_time" bson:"del_time"`                       // 刪除時間
	Balance           uint64  `json:"balance" bson:"balance"`                         // 金額
	TotalBet          uint64  `json:"total_bet" bson:"total_bet"`                     // 總下注
	TotalWin          uint64  `json:"total_win" bson:"total_win"`                     // 總贏分
	TotalEffectiveBet float64 `json:"total_effective_bet" bson:"total_effective_bet"` // 總有效投注
}

type GameWalletOpt struct {
	UserIDs       []uint64
	UserID        uint64
	ClubID        uint64
	TableID       uint64
	DelExpireTime int64
}

func (o *GameWalletOpt) Filter_Mgo() bson.M {
	filter := bson.M{}

	if o == nil {
		return filter
	}

	if len(o.UserIDs) > 0 {
		filter["user_id"] = bson.M{"$in": o.UserIDs}
	} else if o.UserID > 0 {
		filter["user_id"] = o.UserID
	}

	if o.ClubID > 0 {
		filter["club_id"] = o.ClubID
	}

	if o.TableID > 0 {
		filter["table_id"] = o.TableID
	}

	if o.DelExpireTime > 0 {
		filter["del_time"] = bson.M{"$lte": o.DelExpireTime}
	}

	return filter
}
