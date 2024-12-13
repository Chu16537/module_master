package hnats

import "fmt"

const (
	StreamNameGameServer   = "stream_game_server"
	SubjectNameGameServer  = "subject_game_server"
	ConsumerNameGameServer = "consumer_game_server"

	StreamNameGameClientServer   = "stream_game_client_server"
	SubjectNameGameClientServer  = "subject_game_client_server"
	ConsumerNameGameClientServer = "consumer_game_client_server"
)

func subNameRoom(tableID uint64) string {
	return fmt.Sprintf("sub_room_%v", tableID)
}
func consumerNameRoom(tableID uint64) string {
	return fmt.Sprintf("con_room_%v", tableID)
}

func subNamePlayer() string {
	return "sub_player"
}

func consumerRoomToPlayer() string {
	return "con_player"
}
