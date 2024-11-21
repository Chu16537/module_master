package hnats

import (
	"fmt"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mjson"
	"github.com/Chu16537/module_master/mnats"
	"github.com/Chu16537/module_master/proto"
)

// 註冊 table server
func SubTableGameServer(h *mnats.Handler, nodeID int64, subChan chan proto.MQSubData) *errorcode.Error {
	sm := mnats.SubMode{
		Mode: mnats.Sub_Mode_Last,
	}

	err := h.Sub(GameServer, gameserverConsumerKey(nodeID), sm, subChan)
	if err != nil {
		return errorcode.MQSubError(GameServer, err)
	}

	return nil
}

// 推送 gs table > player
func PubTableToPlayer(h *mnats.Handler, p *proto.TableToPlayer) *errorcode.Error {
	data, err := mjson.Marshal(p)
	if err != nil {
		return errorcode.DataMarshalError(fmt.Sprintf("PubTableToPlayer err:%v", err.Error()))
	}

	err = h.Pub(PlayerServer, data)
	if err != nil {
		return errorcode.MQPubError(PlayerServer, data, err)
	}

	return nil
}

// 註冊 player server
func SubPlayerToTable(h *mnats.Handler, tableID uint64, subChan chan proto.MQSubData) *errorcode.Error {
	sm := mnats.SubMode{
		Mode: mnats.Sub_Mode_Last_Ack,
	}

	err := h.Sub(GameServer, tableKey(tableID), sm, subChan)
	if err != nil {
		return errorcode.MQSubError(GameServer, err)
	}

	return nil
}
