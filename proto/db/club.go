package db

import (
	"go.mongodb.org/mongo-driver/bson"
)

const (
	Club_Status_OFF = 0 // 關閉
	Club_Status_ON  = 1 // 開啟
)

type Club struct {
	ID             uint64 `json:"id" bson:"id"`                           // id
	Name           string `json:"name" bson:"name"`                       // 名稱
	PresidentName  string `json:"president_name" bson:"president_name"`   // 會長名稱
	Content        string `json:"content" bson:"content"`                 // 公告
	Status         int    `json:"status" bson:"status"`                   // 狀態
	ExpireTime     int64  `json:"expire_time" bson:"expire_time"`         // 到期時間
	InvitationCode string `json:"invitation_code" bson:"invitation_code"` // 邀請碼 用於創建帳號時使用
}

type ClubOpt struct {
	ClubID         uint64
	ClubIDs        []uint64
	Name           string
	Status         []int
	ExpireTime     int64
	DelExpireTime  int64
	InvitationCode string
	StartTimeUnix  int64
	EndTimeUnix    int64
}

func (o *ClubOpt) Filter_Mgo() bson.M {
	filter := bson.M{}

	if o == nil {
		return filter
	}

	if len(o.ClubIDs) > 0 {
		filter["id"] = bson.M{"$in": o.ClubIDs}
	} else if o.ClubID > 0 {
		filter["id"] = o.ClubID
	}

	if o.Name != "" {
		filter["name"] = o.Name
	}

	if len(o.Status) > 0 {
		filter["status"] = bson.M{"$in": o.Status}
	}

	if o.ExpireTime > 0 {
		filter["expire_time"] = bson.M{"$gte": o.ExpireTime}
	}

	if o.DelExpireTime > 0 {
		filter["expire_time"] = bson.M{"$lte": o.DelExpireTime}
	}

	if o.InvitationCode != "" {
		filter["invitation_code"] = o.InvitationCode
	}

	o.StartTimeUnix, o.EndTimeUnix = checkTimeUxin(o.StartTimeUnix, o.EndTimeUnix)
	if o.StartTimeUnix > 0 {
		if o.EndTimeUnix > 0 {
			filter["expire_time"] = bson.M{"$gte": o.StartTimeUnix, "$lte": o.EndTimeUnix}
		} else {
			filter["expire_time"] = bson.M{"$gte": o.StartTimeUnix}
		}
	}

	return filter
}

type ClubByUserId struct {
	ClubUserInfo ClubUserInfo `bson:"club_user_info"`
	Club         Club         `bson:"club"`
}
