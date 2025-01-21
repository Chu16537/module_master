package proto

type CommonReq struct {
	Platform string `json:"platform"`
	Data     string `json:"data"` // base64 跟 aes加密後資料
}

type CommonRes struct {
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data"`
}

// 查詢共用findopt
type FindOpt struct {
	Start          uint64 `json:"start"`
	Limit          uint64 `json:"limit"`
	StartTimestamp int64  `json:"start_timestamp"`
	EndTimestamp   int64  `json:"end_timestamp"`
}

type LaunchGameReq struct {
	Account string `json:"account"`
	GameId  string `json:"game_id"`
	Lang    string `json:"lang"`
}

type LaunchGameRes struct {
	Url string `json:"url"`
}
