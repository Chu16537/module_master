package db

import (
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// Club_User_Permissions_WaitJoin       = 0 // 申請等待加入
	Club_User_Permissions_Leave          = 0 // 被踢除
	Club_User_Permissions_Member         = 1 // 會員
	Club_User_Permissions_Vice_President = 2 // 副會長
	Club_User_Permissions_President      = 3 // 會長

)

// 會長可設定的權限
var PresidentCanSetPermissions = map[int]struct{}{
	Club_User_Permissions_Leave:          struct{}{},
	Club_User_Permissions_Member:         struct{}{},
	Club_User_Permissions_Vice_President: struct{}{},
}

type ClubUserInfo struct {
	UserID      uint64 `json:"user_id" bson:"user_id"`
	Account     string `json:"account" bson:"account"`
	Password    string `json:"password" bson:"password"`
	NickName    string `json:"nick_name" bson:"nick_name"`
	Token       string `json:"token" bson:"token"`
	ClubID      uint64 `json:"club_id" bson:"club_id"`
	Permissions int    `json:"permissions" bson:"permissions"` // 權限
	Balance     uint64 `json:"balance" bson:"balance"`         // 金額
	TotalBet    uint64 `json:"total_bet" bson:"total_bet"`     // 總下注
	TotalWin    uint64 `json:"total_win" bson:"total_win"`     // 總派彩
}

type ClubUserInfoOpt struct {
	UserID      uint64
	UserIDs     []uint64
	ClubID      uint64
	Account     string
	Password    string
	Permissions []int // 要查詢的權限
	Token       string
	OR          *ClubUserInfoOR
}

type ClubUserInfoOR struct {
	Account  string
	NickName string
}

// 更新權限
type UpdateClubUserInfoPermissions struct {
	UserID      uint64
	ClubID      uint64
	Permissions int
}

func (o *ClubUserInfoOpt) Filter_Mgo() bson.M {
	filter := bson.M{}

	if o == nil {
		return filter
	}

	if len(o.UserIDs) > 0 {
		filter["user_id"] = bson.M{"$in": o.UserIDs}
	} else if o.UserID > 0 {
		filter["user_id"] = o.UserID
	}

	if o.Account != "" {
		filter["account"] = o.Account
	}

	if o.Password != "" {
		filter["password"] = o.Password
	}

	if o.ClubID > 0 {
		filter["club_id"] = o.ClubID
	}

	if len(o.Permissions) > 0 {
		filter["permissions"] = bson.M{"$in": o.Permissions}
	}

	if o.Token != "" {
		filter["token"] = o.Token
	}

	if o.OR != nil {
		or := []bson.M{}

		if o.OR.Account != "" {
			or = append(or, bson.M{"account": o.OR.Account})
		}

		if o.OR.NickName != "" {
			or = append(or, bson.M{"nick_name": o.OR.NickName})
		}

		filter["$or"] = or
	}

	return filter
}
