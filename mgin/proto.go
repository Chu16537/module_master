package mgin

type Req struct {
	Data interface{} `json:"data"`
}

type Res struct {
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data"`
}
