package hnats

import (
	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mnats"
	"github.com/Chu16537/module_master/proto"
)

// 推送 gs
func PubGameServer(h *mnats.Handler, data []byte) *errorcode.Error {
	err := h.Pub(GameServer, data)
	if err != nil {
		return errorcode.MQPubError(GameServer, data, err)
	}

	return nil
}

// 註冊 gs
func SubGameServer(h *mnats.Handler, subChan chan proto.MQSubData) *errorcode.Error {
	sm := mnats.SubMode{
		Mode: mnats.Sub_Mode_Last,
	}

	err := h.Sub(GameServer, sm, subChan)
	if err != nil {
		return errorcode.MQSubError(GameServer, err)
	}

	return nil
}
