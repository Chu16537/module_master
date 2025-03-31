package mwebscoketserver_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/module_master/mwebscoketserver"
)

type TestToken struct {
	Token string `json:"token"`
}

func TestMain(t *testing.T) {
	ctx := context.Background()

	config := &mwebscoketserver.Config{
		Port:               "10000",
		MaxConn:            10000,
		AliveTimeoutSecond: 3,
	}

	a := &aa{}
	err := mwebscoketserver.New(ctx, config, a)
	if err != nil {
		fmt.Println("ws err", err)
		return
	}

	fmt.Println("aa")
	time.Sleep(100 * time.Second)
}

type aa struct{}

func (a *aa) ReadMessage(req *mwebscoketserver.ToHanglerReq) {
	fmt.Println(req.RequestId, req.ClientId, req.Data)

	// b, err := json.Marshal(req.Data)
	// if err != nil {
	// 	fmt.Println("Marshal err", err)
	// }

	s := &TestToken{}
	err := json.Unmarshal(req.Data, s)
	if err != nil {
		fmt.Println("Unmarshal err", err)
	}

	fmt.Println("s", s)
	res := &mwebscoketserver.ToHanglerRes{
		RequestId: req.RequestId,
		ClientId:  req.ClientId,
		Data:      "ss",
	}
	mwebscoketserver.Response(res)
}

func (a *aa) Disconnect(idx uint32) {
	fmt.Println("disconnect", idx)
}
