package mwebscoketserver_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/chu16537/module_master/mwebscoketserver"
)

var (
	userMap map[uint32]mwebscoketserver.IClient
)

type TestToken struct {
	Token string `json:"token"`
}

func TestMain(t *testing.T) {
	userMap = make(map[uint32]mwebscoketserver.IClient)
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

	userMap[client.GetUid()] = client
	return nil
}

func (a *aa) Disconnect(idx uint32) {
	fmt.Println("disconnect", idx)
	delete(userMap, idx)
}

func (a *aa) EventHandler(clientId uint32, req *mwebscoketserver.ClientReq) {

	s := &TestToken{}
	err := json.Unmarshal(req.Data, s)
	if err != nil {
		fmt.Println("Unmarshal err", err)
	}

	res := req.CreateClientRes()
	res.Data = &aa{
		AA: "aa",
		SS: 1,
	}

	resByte, _ := json.Marshal(res)

	userMap[clientId].WriteMessage(resByte)
}
