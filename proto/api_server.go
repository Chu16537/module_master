package proto

// backend_web 格式
type CommRes struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type LoginReq struct {
	Account        string `json:"account"`         // 帳號
	Password       string `json:"password"`        // 密碼
	InvitationCode string `json:"invitation_code"` // 俱樂部邀請碼
}

type LoginRes struct {
	ClubID            uint64 `json:"club_id"`
	ClubName          string `json:"club_name"`
	ClubPresidentName string `json:"club_president_name"`
	ClubContent       string `json:"club_content"`
	UserID            uint64 `json:"user_id"`
	Account           string `json:"account"`
	NickName          string `json:"nick_name"`
	Token             string `json:"token"`
	Permissions       int    `json:"permissions"`
	Balance           uint64 `json:"balance"`
	TotalBet          uint64 `json:"total_bet"`
	TotalWin          uint64 `json:"total_win"`
}

type LogoutReq struct {
}

type LogoutRes struct {
}

type ChangePasswordReq struct {
	NewPassword string `json:"new_password"`
}

type ChangePasswordRes struct {
}

type GetUserDataReq struct {
}

type GetUserDataRes struct {
	Permissions int    `json:"permissions"`
	Balance     uint64 `json:"balance"`
	TotalBet    uint64 `json:"total_bet"`
	TotalWin    uint64 `json:"total_win"`
}

// type TransBalanceReq struct {
// 	TagUserID int64  `json:"tag_user_id"` // 目標玩家
// 	Balance   uint64 `json:"balance"`     // 只可以>0
// }

// type TransBalanceRes struct {
// 	Balance uint64 `json:"balance"` // 轉帳後剩餘金額
// }

// type TransBalanceMainToClubReq struct {
// 	ClubID  uint64 `json:"club_id"` // 俱樂部編號
// 	Balance uint64 `json:"balance"` // 金額
// }

// type TransBalanceMainToClubRes struct {
// 	MainBalance uint64 `json:"main_balance"` // 主錢包
// 	ClubBalance uint64 `json:"club_balance"` // 俱樂部金額
// }
// type TransBalanceClubToRefundReq struct {
// 	ClubID  uint64 `json:"club_id"` // 俱樂部編號
// 	Balance uint64 `json:"balance"` // 金額
// }

// type TransBalanceClubToRefundRes struct {
// 	RefundBalance uint64 `json:"refund_balance"` // 退款錢包
// 	ClubBalance   uint64 `json:"club_balance"`   // 俱樂部金額
// }
// type TransBalanceRefundToMainReq struct {
// 	Balance uint64 `json:"balance"` // 金額
// }

// type TransBalanceRefundToMainRes struct {
// 	MainBalance   uint64 `json:"main_balance"`   // 主錢包
// 	RefundBalance uint64 `json:"refund_balance"` // 退款錢包
// }

// type StartGameReq struct {
// 	GameID int `json:"game_id"`
// }

// type StartGameRes struct {
// 	Url string `json:"url"`
// }

// type GetRecordReq struct {
// 	StartDate string `json:"start_date"`
// 	EndDate   string `json:"end_date"`
// }

// type GetRecordMonthReq struct {
// 	StartDate string `json:"start_date"` // 2024-01-01
// 	EndDate   string `json:"end_date"`   // 2024-03-31 以當天日期為結束
// }

// type GetRecordMonthRes struct {
// 	Record []UserRevenue `json:"record"`
// }

// type UserRevenue struct {
// 	Date    string `json:"date"`    // 2024-03-22
// 	Balance int64  `json:"balance"` // 當天勝負金額
// }

// type GetRecordDayReq struct {
// 	Date string `json:"date"` // 2024-01-01
// }

// type GetRecordDayRes struct {
// 	Record []RecordDay `json:"record"`
// }

// type RecordDay struct {
// 	// 遊戲結果
// }

// type GetTableInfoReq struct {
// }

// type GetTableInfoRes struct {
// }

// type CreateTableReq struct {
// }

// type CreateTableRes struct {
// }

// type UpdateTableReq struct {
// }

// type UpdateTableRes struct {
// }

// type StopTableReq struct {
// }

// type StopTableRes struct {
// }

// type DelTableReq struct {
// }

// type DelTableRes struct {
// }

// type CreateClubReq struct {
// }
// type CreateClubRes struct {
// }
// type DelClubReq struct {
// }
// type DelClubRes struct {
// }

type GetClubInfoReq struct{}

type GetClubInfoRes struct {
	ClubID         uint64 `json:"club_id"`         //  id
	ClubName       string `json:"name"`            // 名稱
	PresidentName  string `json:"president_name"`  // 會長名稱
	Content        string `json:"content"`         // 公告
	InvitationCode string `json:"invitation_code"` // 邀請碼 用於創建帳號時使用
}

type JoinClubReq struct {
	Account        string `json:"account"`         // 帳號
	Password       string `json:"password"`        // 密碼
	NickName       string `json:"nick_name"`       // 暱稱
	InvitationCode string `json:"invitation_code"` // 俱樂部邀請碼
}

type JoinClubRes struct {
}

// type LeaveClubReq struct {
// }
// type LeaveClubRes struct {
// }
// type TransBalanceToClubUserReq struct {
// }
// type TransBalanceToClubUserRes struct {
// }
// type UpdatePermissionsClubReq struct {
// }

// type UpdatePermissionsClubRes struct {
// }
