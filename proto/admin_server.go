package proto

type CreateClubReq struct {
	ClubName   string `json:"club_name"`          // 俱樂部名稱
	Account    string `json:"president_account"`  // 會長帳號
	Password   string `json:"president_password"` // 會長密碼
	Nickname   string `json:"president_nickname"` // 會長名稱
	ExpireTime int64  `json:"expire_time"`        // 到期時間
}

type CreateClubRes struct {
}

type UpdateClubReq struct {
	ClubID     uint64 `json:"club_id"`     // 俱樂部id
	Status     int    `json:"status"`      // 狀態
	ExpireTime int64  `json:"expire_time"` // 到期時間
}

type UpdateClubRes struct {
}

type CreateTableReq struct {
	ClubID     uint64 `json:"club_id"`     // 俱樂部id
	ExpireTime int64  `json:"expire_time"` // 到期時間
	GameID     int    `json:"game_id"`     // 遊戲編號
}

type CreateTableRes struct {
}
