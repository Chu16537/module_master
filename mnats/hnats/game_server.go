package hnats

import (
	"github.com/Chu16537/module_master/mjson"
	"github.com/Chu16537/module_master/mnats"
	"github.com/Chu16537/module_master/proto"
)

// game server 註冊
func SubGameServer(h *mnats.Handler, subChan chan proto.MQSubData) error {
	err := h.CreateSubjects(StreamNameGameServer, SubjectNameGameServer)
	if err != nil {
		return err
	}

	sm := mnats.SubMode{
		Mode: mnats.Sub_Mode_Last,
	}

	err = h.Sub(SubjectNameGameServer, ConsumerNameGameServer, sm, subChan)
	if err != nil {
		return err
	}

	return nil
}

// 玩家 > 房間 註冊
func SubRoom(h *mnats.Handler, tableID uint64, subChan chan proto.MQSubData) error {
	sm := mnats.SubMode{
		Mode: mnats.Sub_Mode_Last_Ack,
	}

	subName := subNameRoom(tableID)
	conName := consumerNameRoom(tableID)

	err := h.CreateSubjects(StreamNameGameServer, subName)
	if err != nil {
		return err
	}

	err = h.Sub(subName, conName, sm, subChan)
	if err != nil {
		return err
	}

	return nil
}

// 刪除房間所有訊息
func DelRoom(h *mnats.Handler, tableID uint64) error {
	subName := subNameRoom(tableID)
	err := h.DelSubject(StreamNameGameServer, subName)
	if err != nil {
		return err
	}

	return nil
}

// 房間推播給玩家
func PubPlayer(h *mnats.Handler, data *proto.TableToPlayer) error {

	b, err := mjson.Marshal(data)
	if err != nil {
		return err
	}

	err = h.Pub(subNamePlayer(), b)
	if err != nil {
		return err
	}

	return nil
}
