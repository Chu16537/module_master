package mwebscoketserver

import "encoding/json"

type ClientReq struct {
	RequestId string          `json:"req_id"` // 前端創建的req id
	Command   string          `json:"command"`
	Data      json.RawMessage `json:"data"`
}

func (c *ClientReq) CreateClientRes() *ClientRes {
	return &ClientRes{
		RequestId: c.RequestId,
		Command:   c.Command,
	}
}

type ClientRes struct {
	RequestId string      `json:"req_id"` // 前端創建的req id
	Command   string      `json:"command"`
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data"`
}
