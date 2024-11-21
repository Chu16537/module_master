package hnats

import "fmt"

const (
	GameServer   = "game_server"
	PlayerServer = "player_server"
	Table        = "table"
)

func gameserverConsumerKey(nodeId int64) string {
	return fmt.Sprintf("%v.%v", GameServer, nodeId)
}

func tableKey(tableId uint64) string {
	return fmt.Sprintf("%v.%v", Table, tableId)
}
