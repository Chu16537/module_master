package proto

import "github.com/Chu16537/module_master/proto/db"

type CommonReq struct {
	Platform string `json:"platform"`
	Data     string `json:"data"` // base64 跟 aes加密後資料
}

type CommonRes struct {
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data"`
}

// 查詢共用findopt
type FindOpt struct {
	Start         uint64 `json:"start"`
	Limit         uint64 `json:"limit"`
	StartTimeUnix int64  `json:"start_timeunix"`
	EndTimeUnix   int64  `json:"end_timeunix"`
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

type TransBalanceClubeReq struct {
	TagUserID uint64 `json:"tagUser_id"`
	Amount    uint64 `json:"amount"`
}
type TransBalanceClubeRes struct {
	Balance uint64 `json:"balance"`
}

type UpdateClubContentReq struct {
	Content string `json:"content"`
}
type UpdateClubContentRes struct {
}

type UpdateClubMemberPermissionsReq struct {
	TagUserID   uint64 `json:"tagUser_id"`
	Permissions int    `json:"permissions"`
}

type UpdateClubMemberPermissionsRes struct {
}

type GetTableReq struct {
	FindOpt
	ClubID uint64 `json:"club_id"`
	Status []int  `json:"status"`
	GameID int64  `json:"game_id"`
}

type GetTableRes struct {
	Tables []*db.Table `json:"tables"`
}

type UpdateTableGameReq struct {
	TableID    uint64 `json:"table_id"`
	GameConfig []byte `json:"game_config"` // 每個遊戲設定不同
}

type UpdateTableGameRes struct {
	Table *db.Table `json:"table"`
}

type UpdateTableReq struct {
	TableID    uint64 `json:"table_id"`
	Status     int    `json:"status"`
	ExpireTime int64  `json:"expire_time"`
}

type UpdateTableRes struct {
}

type GetUserRecordTotalReq struct {
	FindOpt
}

type GetUserRecordTotalRes struct {
	RecordTotalData []RecordTotalData `json:"datas"`
}

type RecordTotalData struct {
	Date  string `json:"date"`
	Total int64  `json:"total"`
}

type GetUserRecordReq struct {
	FindOpt
}

type GetUserRecordRes struct {
	UserRecordData []UserRecordData `json:"datas"`
	Total          int64            `json:"total"`
}

type UserRecordData struct {
	UserID        uint64 `json:"user_id"`
	GameRecordID  string `json:"game_record_id"`
	CreateTime    int64  `json:"create_time"`
	ClubID        uint64 `json:"club_id"`
	TableID       uint64 `json:"table_id"`
	GameID        int    `json:"game_id"`
	GameType      int    `json:"game_type"`
	ResultBalance int64  `json:"result_balance"`
	GameResult    []byte `json:"game_result"`
	Info          []byte `json:"info"`
}

type LaunchGameReq struct {
}

type LaunchGameRes struct {
	Url string `json:"url"`
}
