package tableserverclient

import (
	"context"
	"fmt"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mgrpc/commongrpc"
	"github.com/Chu16537/module_master/mjson"
	"github.com/Chu16537/module_master/proto/db"
)

// 取得牌桌
func (h *Handler) GetTable(ctx context.Context, logTracer string, tableOpt *db.TableOpt, findOpt *db.FindOpt) ([]*db.Table, *errorcode.Error) {
	reqData := &TSGetTableReq{
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

	resData := &TSGetTableRes{}
	err = mjson.Unmarshal(res.Data, resData)
	if err != nil {
		return nil, errorcode.New(errorcode.Code_Data_Unmarshal_Error, fmt.Errorf("GetTable Unmarshal error:%v", err))
	}

	return resData.Tables, nil
}

// 更新牌桌 遊戲設定
func (h *Handler) UpdateTableGame(ctx context.Context, logTracer string, tableOpt *db.TableOpt, gameConfig []byte) (*db.Table, *errorcode.Error) {
	reqData := &TSUpdateTableGameReq{
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

	resData := &TSUpdateTableGameRes{}
	err = mjson.Unmarshal(res.Data, resData)
	if err != nil {
		return nil, errorcode.New(errorcode.Code_Data_Unmarshal_Error, fmt.Errorf("UpdateTableGame Unmarshal error:%v", err))
	}

	return resData.Table, nil
}

// 更新牌桌狀態
func (h *Handler) UpdateTable(ctx context.Context, logTracer string, tableOpt *db.TableOpt, status int, expireTime int64) *errorcode.Error {
	reqData := &TSUpdateTableReq{
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

	// resData := &TSUpdateTableRes{}
	// err = mjson.Unmarshal(res.Data, resData)
	// if err != nil {
	// 	return errorcode.New(errorcode.Code_Data_Unmarshal_Error,fmt.Errorf("UpdateTableGame Unmarshal error:%v", err))
	// }

	return nil
}

// 創建牌桌
func (h *Handler) CreateTable(ctx context.Context, logTracer string, clubID uint64, expireTime int64, gameID int) *errorcode.Error {
	reqData := &TSCreateTableReq{
		ClubID:     clubID,
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

	// resData := &TSCreateTableRes{}
	// err = mjson.Unmarshal(res.Data, resData)
	// if err != nil {
	// 	return errorcode.New(errorcode.Code_Data_Unmarshal_Error,fmt.Errorf("CreateTable Unmarshal error:%v", err))
	// }

	return nil
}
