package cache

type User struct {
	UserID      uint64 `json:"user_id"`
	ClubID      uint64 `json:"club_id"`
	NickName    string `json:"nick_name"`
	Permissions int    `json:"permissions"`
}
