package proto

type MQSubData struct {
	SequenceID uint64 `json:"sequence_id"`
	Data       []byte `json:"data"`
}

// table server 刪除房間
type TableServerDeleteRoom struct {
	TableIDs []uint64 // 要刪除的房間編號
}

// 玩家請求 推給房間
type PlayerReqPushToRoom struct {
	ReqID     string      `json:"req_id"`     // 請求編號
	UserID    uint64      `json:"user_id"`    // 玩家編號
	TableName uint64      `json:"table_name"` // 房間編號
	EventCode uint64      `json:"event_code"` // 事件編號
	Data      interface{} `json:"data"`       // 資料
}

// 房間推送資料給玩家
type RoomPushToUser struct {
	ReqID     string      `json:"req_id"`     // 請求編號
	UserIDs   []uint64    `json:"user_ids"`   // 玩家編號
	RoomName  uint64      `json:"room_name"`  // 房間編號
	EventCode uint64      `json:"event_code"` // 事件編號
	Data      interface{} `json:"data"`       // 資料
	ErrorCode int         `json:"error_code"` // 錯誤碼
}
