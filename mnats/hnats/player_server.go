package hnats

// // 推送 player > gs table
// func PubPlayerToTable(h *mnats.Handler, p *proto.PlayerToTable) *errorcode.Error {
// 	data, err := mjson.Marshal(p)
// 	if err != nil {
// 		return errorcode.DataMarshalError(fmt.Sprintf("PubPlayerToTable err:%v", err.Error()))
// 	}

// 	err = h.Pub(tableKey(p.TableID), data)
// 	if err != nil {
// 		return errorcode.MQPubError(GameServer, data, err)
// 	}

// 	return nil
// }

// // 註冊 game server
// func SubTableToPlayer(h *mnats.Handler, nodeID int64, subChan chan proto.MQSubData) *errorcode.Error {
// 	sm := mnats.SubMode{
// 		Mode: mnats.Sub_Mode_Last,
// 	}

// 	err := h.Sub(PlayerServer, PlayerServer, gameserverConsumerKey(nodeID), sm, subChan)
// 	if err != nil {
// 		return errorcode.MQSubError(PlayerServer, err)
// 	}

// 	return nil
// }
