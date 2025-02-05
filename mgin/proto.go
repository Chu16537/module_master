package mgin

type Req struct {
	Platform string `json:"platform"`
	Data     string `json:"data"` // base64 跟 aes加密後資料
}

type Res struct {
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data"`
}
