package errorcode

const (
	SuccessCode = 0 // 成功
	ErrorCode   = 1 // 錯誤

	// 系統
	Server_Error         = 10001 // server錯誤
	Data_Unmarshal_Error = 10002 // 資料Unmarshal錯誤
	Data_Marshal_Error   = 10003 // 資料Unmarshal錯誤
	Data_Error           = 10004 // 資料錯誤 Unmarshal 成功但內容不正確
	Data_Is_Exist        = 10005 // 資料已經存在
	Data_Not_Exist       = 10006 // 資料不存在
	Timeout              = 10007 // 超時
	Token_Not_Exist      = 10008 // token不存在
	Handler_Not_Exist    = 10009 // handler 不存在

	// 俱樂部
	Club_Not_Exist            = 20001 // 俱樂部不存在
	Club_User_Not_Permissions = 20002 // 權限不足
	Club_User_Not_In_Club     = 20003 // 會員不在俱樂部
	Club_User_Balacne_Less    = 20004 // 玩家金額不足
	Club_User_Account_Exist   = 20005 // 玩家帳號已經存在

	// 牌桌
	Table_Not_Exist       = 30001  // 牌桌不存在
	Table_LV_Insufficient = 30002  // 牌桌等級不足
	Table_Status_Error    = 30003  // 牌桌狀態錯誤
	Game_Not_Type         = 30004  // 沒有此遊戲類型
	Game_ID_Error         = 30005  // 遊戲ID錯誤
	Game_Setting_Error    = 30006  // 遊戲設定驗證錯誤
	Game_Setting_Error_1  = 300061 // 遊戲設定驗證錯誤
	Game_Setting_Error_2  = 300062 // 遊戲設定驗證錯誤
	Game_Setting_Error_3  = 300063 // 遊戲設定驗證錯誤
	Game_Setting_Error_4  = 300064 // 遊戲設定驗證錯誤
	Game_Setting_Error_5  = 300065 // 遊戲設定驗證錯誤

	// 遊戲
	GameWallet_Not_Exist    = 40001 // 錢包不存在
	GameWallet_Balacne_Less = 40002 // 錢包金額不足

	// Redis

	// api_server
	Account_Length_Error   = 60001 // account 長度錯誤
	Password_Length_Error  = 60002 // password 長度錯誤
	NickName_Length_Error  = 60003 // nickname 長度錯誤
	Account_Password_Error = 60003 // 帳號或密碼錯誤

	// 內部服務
	RoomServer_Req_Error = 70001 // room server 請求錯誤
)
