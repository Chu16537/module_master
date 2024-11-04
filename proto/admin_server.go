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

type CreateTableReq struct {
	ClubID     string `json:"club_id"`     // 俱樂部名稱id
	ExpireTime int64  `json:"expire_time"` // 到期時間
}

type CreateTableRes struct {
}
