package mwebscoketserver

import (
	"fmt"
	"strconv"
	"strings"
)

// 玩家請求
type ClientReq struct {
	Id   string      `json:"id"` // 前端創建的req id
	Data interface{} `json:"data"`
}

func (r *ClientReq) NewID(clientId int64) {
	r.Id = fmt.Sprintf("%v_%v", r.Id, clientId)
}

func (r *ClientReq) GetId() string {
	return r.Id
}

// 玩家回覆
type ClientRes struct {
	Id   string      `json:"id"`
	Data interface{} `json:"data"`
}

func (r *ClientRes) GetId() (string, int64) {
	parts := strings.Split(r.Id, "_")
	// 取得第一個 前端創建的id
	reqId := parts[0]

	// 取得最後一個 clientId
	clientIdStr := parts[len(parts)-1]
	clientId, _ := strconv.Atoi(clientIdStr)

	return reqId, int64(clientId)
}
