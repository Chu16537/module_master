package test_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Chu16537/module_master/mgrpc/commongrpc"
	"github.com/Chu16537/module_master/mgrpc/mgrpcclient"
	"github.com/Chu16537/module_master/mjson"
	"github.com/Chu16537/module_master/proto"
	"github.com/Chu16537/module_master/proto/db"
)

func TestClient(t *testing.T) {
	config := &mgrpcclient.Config{
		Addr:          ":50051",
		TimeoutSecond: 3,
	}

	h, err := mgrpcclient.New(context.Background(), config)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	reqData := &db.TableOpt{
		Status: []int{0},
	}

	rb, err := mjson.Marshal(reqData)
	if err != nil {
		fmt.Println("err 2", err)
		return
	}

	req := &commongrpc.UnaryRPCReq{
		LogData: &commongrpc.LogData{
			Tracer: "test_client",
		},
		EventCode: proto.RS_GET_TABLE,
		Data:      rb,
	}
	res, errC := h.UnaryRPC(context.Background(), req)
	if errC != nil {
		fmt.Println("errC", errC)
	}

	fmt.Println("res", res)

	res, errC = h.UnaryRPC(context.Background(), req)
	if errC != nil {
		fmt.Println("errC", errC)
	}

	fmt.Println("res", res)

	res, errC = h.UnaryRPC(context.Background(), req)
	if errC != nil {
		fmt.Println("errC", errC)
	}

	fmt.Println("res", res)

}
