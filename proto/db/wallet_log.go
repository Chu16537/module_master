package db

import (
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// 錢包類型
	WalletType_Club = 1 // 俱樂部錢包
	WalletType_Game = 2 // 遊戲錢包

	// 事件類型
	EventType_UserTransBalanceIn         = 1 // 玩家 轉入給 玩家
	EventType_UserTransBalanceOut        = 2 // 玩家 轉出給 玩家
	EventType_UserTransBalanceClubToGame = 3 // 俱樂部 轉入給 遊戲
	EventType_UserTransBalanceGameToClub = 4 // 遊戲 轉入給 俱樂部
	EventType_GameBet                    = 5 // 遊戲下注
	EventType_GamePay                    = 6 // 遊戲派彩

)

// 錢包log
type WalletLog struct {
	OrderID       string `json:"order_id" bson:"order_id"`             // 注單唯一碼
	UserID        uint64 `json:"user_id" bson:"user_id"`               // 玩家id
	WalletType    int    `json:"wallet_type" bson:"wallet_type"`       // 錢包類型 (ex 主錢包 遊戲錢包 ... 等)
	EventType     int    `json:"event_type" bson:"event_type"`         // 動作類型 (ex 進出房間 遊戲贏分 ... 等)
	UpdateBalance int64  `json:"update_balance" bson:"update_balance"` // 修改金額
	BeforeBalance uint64 `json:"before_balance" bson:"before_balance"` // 之前金額
	Afterbalance  uint64 `json:"after_balance" bson:"after_balance"`   // 之後金額
	CreateTime    int64  `json:"create_time" bson:"create_time"`       // 創建時間
	ClubID        uint64 `json:"club_id" bson:"club_id"`               // club id
	TableID       uint64 `json:"table_id" bson:"table_id"`             // table id
	CustomInfo    string `json:"custom_info" bson:"custom_info"`       // 客製化訊息(json)
}

type WalletLogOpt struct {
	OrderID       string
	UserID        uint64
	UserIDs       []uint64
	ClubID        uint64
	ClubIDs       []uint64
	EventType     int
	WalletType    int
	StartTimeUnix int64
	EndTimeUnix   int64
}

func (o *WalletLogOpt) Filter_Mgo() bson.M {
	filter := bson.M{}

	if o == nil {
		return filter
	}

	if o.OrderID != "" {
		filter["order_id"] = o.OrderID
	}

	if len(o.UserIDs) > 0 {
		filter["id"] = bson.M{"$in": o.UserIDs}
	} else if o.UserID > 0 {
		filter["id"] = o.UserID
	}

	if len(o.ClubIDs) > 0 {
		filter["id"] = bson.M{"$in": o.ClubIDs}
	} else if o.ClubID > 0 {
		filter["id"] = o.ClubID
	}

	if o.EventType > 0 {
		filter["event_type"] = o.EventType
	}

	if o.WalletType > 0 {
		filter["wallet_type"] = o.WalletType
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

type WalletInfoClub struct {
	OrderID              string
	UserID               uint64
	UserAmount           int64
	UserBeforeBalance    uint64
	TagUsertID           uint64
	TagUserAmount        int64
	TagUserBeforeBalance uint64
	CreateTime           int64
	ClubID               uint64
}

// 轉帳 log Club
func (w *WalletInfoClub) Log() []*WalletLog {
	user := &WalletLog{
		OrderID:       w.OrderID,
		UserID:        w.UserID,
		EventType:     EventType_UserTransBalanceOut,
		WalletType:    WalletType_Club,
		UpdateBalance: w.UserAmount,
		BeforeBalance: w.UserBeforeBalance,
		Afterbalance:  w.UserBeforeBalance + uint64(w.UserAmount),
		CreateTime:    w.CreateTime,
		ClubID:        w.ClubID,
	}

	tagUser := &WalletLog{
		OrderID:       w.OrderID,
		UserID:        w.TagUsertID,
		EventType:     EventType_UserTransBalanceIn,
		WalletType:    WalletType_Club,
		UpdateBalance: w.TagUserAmount,
		BeforeBalance: w.TagUserBeforeBalance,
		Afterbalance:  w.TagUserBeforeBalance + uint64(w.TagUserAmount),
		CreateTime:    w.CreateTime,
		ClubID:        w.ClubID,
	}

	return []*WalletLog{user, tagUser}
}

type WalletInfoClubToGame struct {
	OrderID           string
	UserID            uint64
	UpdateGameAmount  int64
	GameBeforeBalance uint64
	UpdateClubAmount  int64
	ClubBeforeBalance uint64
	CreateTime        int64
	ClubID            uint64
	TableID           uint64
}

// 俱樂部 > 遊戲
func (w *WalletInfoClubToGame) Log() []*WalletLog {
	l1 := &WalletLog{
		OrderID:       w.OrderID,
		UserID:        w.UserID,
		EventType:     EventType_UserTransBalanceClubToGame,
		WalletType:    WalletType_Club,
		UpdateBalance: w.UpdateClubAmount,
		BeforeBalance: w.ClubBeforeBalance,
		Afterbalance:  w.ClubBeforeBalance + uint64(w.UpdateClubAmount),
		CreateTime:    w.CreateTime,
		ClubID:        w.ClubID,
		TableID:       w.TableID,
	}

	l2 := &WalletLog{
		OrderID:       w.OrderID,
		UserID:        w.UserID,
		EventType:     EventType_UserTransBalanceClubToGame,
		WalletType:    WalletType_Game,
		UpdateBalance: w.UpdateGameAmount,
		BeforeBalance: w.GameBeforeBalance,
		Afterbalance:  w.GameBeforeBalance + uint64(w.UpdateGameAmount),
		CreateTime:    w.CreateTime,
		ClubID:        w.ClubID,
		TableID:       w.TableID,
	}

	return []*WalletLog{l1, l2}
}

type WalletInfoGameToClub struct {
	OrderID           string
	UserID            uint64
	UpdateGameAmount  int64
	GameBeforeBalance uint64
	UpdateClubAmount  int64
	ClubBeforeBalance uint64
	CreateTime        int64
	ClubID            uint64
	TableID           uint64
}

// 遊戲 > 俱樂部
func (w *WalletInfoGameToClub) Log() []*WalletLog {
	l1 := &WalletLog{
		OrderID:       w.OrderID,
		UserID:        w.UserID,
		EventType:     EventType_UserTransBalanceGameToClub,
		WalletType:    WalletType_Club,
		UpdateBalance: w.UpdateClubAmount,
		BeforeBalance: w.ClubBeforeBalance,
		Afterbalance:  w.ClubBeforeBalance + uint64(w.UpdateClubAmount),
		CreateTime:    w.CreateTime,
		ClubID:        w.ClubID,
		TableID:       w.TableID,
	}

	l2 := &WalletLog{
		OrderID:       w.OrderID,
		UserID:        w.UserID,
		EventType:     EventType_UserTransBalanceGameToClub,
		WalletType:    WalletType_Game,
		UpdateBalance: w.UpdateGameAmount,
		BeforeBalance: w.GameBeforeBalance,
		Afterbalance:  w.GameBeforeBalance + uint64(w.UpdateGameAmount),
		CreateTime:    w.CreateTime,
		ClubID:        w.ClubID,
		TableID:       w.TableID,
	}

	return []*WalletLog{l1, l2}
}

// 遊戲結果 下注 跟 派獎
type WalletInfoGameRecode struct {
	OrderID       string
	UserID        uint64
	Bet           uint64
	Win           uint64
	BeforeBalance uint64
	CreateTime    int64
	ClubID        uint64
	TableID       uint64
}

// 下注 跟 派獎
func (w *WalletInfoGameRecode) Log() []*WalletLog {
	l1 := &WalletLog{
		OrderID:       w.OrderID,
		UserID:        w.UserID,
		EventType:     EventType_GameBet,
		WalletType:    WalletType_Game,
		UpdateBalance: int64(-w.Bet),
		BeforeBalance: w.BeforeBalance,
		Afterbalance:  w.BeforeBalance - w.Bet,
		CreateTime:    w.CreateTime,
		ClubID:        w.ClubID,
		TableID:       w.TableID,
	}

	l2 := &WalletLog{
		OrderID:       w.OrderID,
		UserID:        w.UserID,
		EventType:     EventType_GamePay,
		WalletType:    WalletType_Game,
		UpdateBalance: int64(w.Win),
		BeforeBalance: w.BeforeBalance - w.Bet,
		Afterbalance:  w.BeforeBalance - w.Bet + w.Win,
		CreateTime:    w.CreateTime,
		ClubID:        w.ClubID,
		TableID:       w.TableID,
	}

	return []*WalletLog{l1, l2}
}
