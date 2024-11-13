package proto

const (
	GS_CREATE_ROOM = 1
)

type GSCreateRoomReq struct {
	TableID uint64 `json:"tables"`
}

type GSCreateRoomRes struct {
}
