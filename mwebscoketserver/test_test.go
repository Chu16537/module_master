package mwebscoketserver_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/module_master/mwebscoketserver"
)

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
