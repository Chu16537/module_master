package mwebscoketserver

type ToHanglerReq struct {
	RequestId string      `json:"req_id"` // 前端創建的req id
	ClientId  uint32      `json:"client_id"`
	Data      interface{} `json:"data"`
}

type ToHanglerRes struct {
	RequestId string      `json:"req_id"` // 前端創建的req id
	ClientId  uint32      `json:"client_id"`
	Data      interface{} `json:"data"`
}

type ClientRes struct {
	RequestId string      `json:"req_id"` // 前端創建的req id
	Data      interface{} `json:"data"`
}
