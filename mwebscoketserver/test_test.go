package mwebscoketserver_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/module_master/mlog"
	"github.com/Chu16537/module_master/muid"
	"github.com/Chu16537/module_master/mwebscoketserver"
)

func TestMain(t *testing.T) {
	ctx := context.Background()

	config := &mwebscoketserver.Config{
		Addr:               "0.0.0.0:10000",
		MaxConn:            10000,
		AliveTimeoutSecond: 10,
	}

	uid := muid.New(1)

	logConfig := &mlog.Config{
		Name: "test_ws",
	}
	err := mlog.New(logConfig)
	if err != nil {
		fmt.Println("log", err)
		return
	}

	l := mlog.Get("test_ws")

	a := &aa{}
	err = mwebscoketserver.New(ctx, config, uid, l, a)
	if err != nil {
		fmt.Println("ws err", err)
		return
	}

	fmt.Println("aa")
	time.Sleep(100 * time.Second)
}

type aa struct{}

func (a *aa) ReadMessage(req *mwebscoketserver.ClientReq) {

	fmt.Println(req.Id)

	mwebscoketserver.Response(nil)
}
