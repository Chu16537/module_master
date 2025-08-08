package tableserverclient

import (
	"context"
	"fmt"

	"github.com/chu16537/module_master/errorcode"
	"github.com/chu16537/module_master/mgrpc/commongrpc"
	"github.com/chu16537/module_master/mjson"
	"github.com/chu16537/module_master/proto/db"
)

// 取得牌桌
func GetTable(ctx context.Context, logTracer string, tableOpt *db.TableOpt, findOpt *db.FindOpt) ([]*db.Table, *errorcode.Error) {
	if h == nil {
		return nil, errorcode.New(errorcode.Code_Server_Error, fmt.Errorf("client nil"))
	}

	reqData := &GetTableReq{
		TableOpt: tableOpt,
		FindOpt:  findOpt,
	}

	reqDataBytes, err := mjson.Marshal(reqData)
	if err != nil {
		return nil, errorcode.New(errorcode.Code_Data_Marshal_Error, fmt.Errorf("GetTable Marshal error:%v", err))
	}

	req := &commongrpc.UnaryRPCReq{
		LogData: &commongrpc.LogData{
			Tracer: logTracer,
		},
		EventCode: GET_TABLE,
		Data:      reqDataBytes,
	}

	res, errC := h.client.UnaryRPC(ctx, req)
	if errC != nil {
		return nil, errC
	}

	resData := &GetTableRes{}
	err = mjson.Unmarshal(res.Data, resData)
	if err != nil {
		return nil, errorcode.New(errorcode.Code_Data_Unmarshal_Error, fmt.Errorf("GetTable Unmarshal error:%v", err))
	}

	return resData.Tables, nil
}

// 更新牌桌 遊戲設定
func UpdateTableGame(ctx context.Context, logTracer string, tableOpt *db.TableOpt, gameConfig []byte) (*db.Table, *errorcode.Error) {
	if h == nil {
		return nil, errorcode.New(errorcode.Code_Server_Error, fmt.Errorf("client nil"))
	}

	reqData := &UpdateTableGameReq{
		TableOpt:   tableOpt,
		GameConfig: gameConfig,
	}

	reqDataBytes, err := mjson.Marshal(reqData)
	if err != nil {
		return nil, errorcode.New(errorcode.Code_Data_Marshal_Error, fmt.Errorf("UpdateTableGame Marshal error:%v", err))
	}

	req := &commongrpc.UnaryRPCReq{
		LogData: &commongrpc.LogData{
			Tracer: logTracer,
		},
		EventCode: UPDATE_TABLE_GAME,
		Data:      reqDataBytes,
	}

	res, errC := h.client.UnaryRPC(ctx, req)
	if errC != nil {
		return nil, errC
	}

	resData := &UpdateTableGameRes{}
	err = mjson.Unmarshal(res.Data, resData)
	if err != nil {
		return nil, errorcode.New(errorcode.Code_Data_Unmarshal_Error, fmt.Errorf("UpdateTableGame Unmarshal error:%v", err))
	}

	return resData.Table, nil
}

// 更新牌桌狀態
func UpdateTable(ctx context.Context, logTracer string, tableOpt *db.TableOpt, status int, expireTime int64) *errorcode.Error {
	if h == nil {
		return errorcode.New(errorcode.Code_Server_Error, fmt.Errorf("client nil"))
	}

	reqData := &UpdateTableReq{
		TableOpt:   tableOpt,
		Status:     status,
		ExpireTime: expireTime,
	}

	reqDataBytes, err := mjson.Marshal(reqData)
	if err != nil {
		return errorcode.New(errorcode.Code_Data_Marshal_Error, fmt.Errorf("UpdateTableStatus Marshal error:%v", err))
	}

	req := &commongrpc.UnaryRPCReq{
		LogData: &commongrpc.LogData{
			Tracer: logTracer,
		},
		EventCode: UPDATE_TABLE_STATUS,
		Data:      reqDataBytes,
	}

	_, errC := h.client.UnaryRPC(ctx, req)
	if errC != nil {
		return errC
	}

	// resData := &UpdateTableRes{}
	// err = mjson.Unmarshal(res.Data, resData)
	// if err != nil {
	// 	return errorcode.New(errorcode.Code_Data_Unmarshal_Error,fmt.Errorf("UpdateTableGame Unmarshal error:%v", err))
	// }

	return nil
}

// 創建牌桌
func CreateTable(ctx context.Context, logTracer string, platformID uint64, expireTime int64, gameID int) *errorcode.Error {
	if h == nil {
		return errorcode.New(errorcode.Code_Server_Error, fmt.Errorf("client nil"))
	}

	reqData := &CreateTableReq{
		PlatformID: platformID,
		ExpireTime: expireTime,
		GameID:     gameID,
	}

	reqDataBytes, err := mjson.Marshal(reqData)
	if err != nil {
		return errorcode.New(errorcode.Code_Data_Marshal_Error, fmt.Errorf("CreateTable Marshal error:%v", err))
	}

	req := &commongrpc.UnaryRPCReq{
		LogData: &commongrpc.LogData{
			Tracer: logTracer,
		},
		EventCode: CREATE_TABLE,
		Data:      reqDataBytes,
	}

	_, errC := h.client.UnaryRPC(ctx, req)
	if errC != nil {
		return errC
	}

	// resData := &CreateTableRes{}
	// err = mjson.Unmarshal(res.Data, resData)
	// if err != nil {
	// 	return errorcode.New(errorcode.Code_Data_Unmarshal_Error,fmt.Errorf("CreateTable Unmarshal error:%v", err))
	// }

	return nil
}
