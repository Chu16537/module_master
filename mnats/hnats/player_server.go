package hnats

import (
	"github.com/Chu16537/module_master/mjson"
	"github.com/Chu16537/module_master/mnats"
	"github.com/Chu16537/module_master/proto"
)

func SubGameClientServer(h *mnats.Handler, subChan chan proto.MQSubData) error {
	sm := mnats.SubMode{
		Mode: mnats.Sub_Mode_Last,
	}

	err := h.CreateSubjects(StreamNameGameClientServer, SubjectNameGameClientServer)
	if err != nil {
		return err
	}

	err = h.Sub(SubjectNameGameClientServer, ConsumerNameGameClientServer, sm, subChan)
	if err != nil {
		return err
	}

	return nil
}

// room > player
func SubPlayer(h *mnats.Handler, subChan chan proto.MQSubData) error {
	subName := subNamePlayer()
	conName := consumerRoomToPlayer()

	sm := mnats.SubMode{
		Mode: mnats.Sub_Mode_Last,
	}

	err := h.CreateSubjects(StreamNameGameClientServer, subName)
	if err != nil {
		return err
	}

	err = h.Sub(subName, conName, sm, subChan)
	if err != nil {
		return err
	}

	return nil
}

// 房間推播給玩家
func PubRoom(h *mnats.Handler, tableID uint64, data *proto.PlayerToTable) error {
	b, err := mjson.Marshal(data)
	if err != nil {
		return err
	}

	err = h.Pub(subNameRoom(tableID), b)
	if err != nil {
		return err
	}

	return nil
}
