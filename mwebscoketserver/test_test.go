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
		MaxConn:            1,
		AliveTimeoutSecond: 5,
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

type aa struct {
	AA string `json:"aa"`
	SS int    `json:"ss"`
}

func (a *aa) Connect(client mwebscoketserver.IClient) error {
	fmt.Println("Connect", client.GetUid())
	return nil
}

func (a *aa) Disconnect(idx uint32) {
	fmt.Println("disconnect", idx)
}

func (a *aa) ReadMessage(toHanglerReq *mwebscoketserver.ToHanglerReq) {
	// fmt.Println(toHanglerReq.Req.RequestId, toHanglerReq.ClientId, toHanglerReq.Req.Data)

	s := &TestToken{}
	err := json.Unmarshal(toHanglerReq.Req.Data, s)
	if err != nil {
		fmt.Println("Unmarshal err", err)
	}

	res := toHanglerReq.Req.CreateClientRes()
	res.Data = &aa{
		AA: "aa",
		SS: 1,
	}
	toHanglerRes := &mwebscoketserver.ToHanglerRes{
		ClientId: toHanglerReq.ClientId,
		Res:      res,
	}
	mwebscoketserver.Response(toHanglerRes)
}
