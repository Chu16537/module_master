package tableserverclient

import (
	"context"
	"fmt"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mgrpc/commongrpc"
	"github.com/Chu16537/module_master/mjson"
)

// gs 斷線
func (h *handler) GameServerDone(ctx context.Context, logTracer string, nodeID int64) (bool, *errorcode.Error) {
	reqData := &GameServerDoneReq{
		NodeID: nodeID,
	}

	reqDataBytes, err := mjson.Marshal(reqData)
	if err != nil {
		return false, errorcode.New(errorcode.Code_Data_Marshal_Error, fmt.Errorf("GameServerDone Marshal error:%v", err))
	}

	req := &commongrpc.UnaryRPCReq{
		LogData: &commongrpc.LogData{
			Tracer: logTracer,
		},
		EventCode: GET_TABLE,
		Data:      reqDataBytes,
	}

	_, errC := h.client.UnaryRPC(ctx, req)
	if errC != nil {
		return false, errC
	}

	// resData := &GameServerDoneRes{}
	// err = mjson.Unmarshal(res.Data, resData)
	// if err != nil {
	// 	return false, errorcode.DataUnmarshalError(fmt.Sprintf("GameServerDone Unmarshal error:%v", err))
	// }

	return true, nil

}
