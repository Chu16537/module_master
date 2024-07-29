package proto

type PubToRoom struct {
	ReqID     string      `json:"req_id"`     // 請求編號
	UserID    uint64      `json:"user_id"`    // 玩家編號
	RoomName  uint64      `json:"room_name"`  // 房間編號
	EventCode uint64      `json:"event_code"` // 事件編號
	Data      interface{} `json:"data"`       // 資料
}

type PubToGame struct {
	ReqID     string      `json:"req_id"`     // 請求編號
	UserIDs   []uint64    `json:"user_ids"`   // 玩家編號
	RoomName  uint64      `json:"room_name"`  // 房間編號
	EventCode uint64      `json:"event_code"` // 事件編號
	Data      interface{} `json:"data"`       // 資料
	ErrorCode int         `json:"error_code"` // 錯誤碼
}
