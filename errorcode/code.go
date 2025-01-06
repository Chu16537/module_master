package errorcode

const (
	Code_Success = 0 // 成功
	Code_Error   = 1 // 錯誤

	// 系統
	Code_Server_Error         = 10001 // server錯誤
	Code_Data_Unmarshal_Error = 10002 // 資料Unmarshal錯誤
	Code_Data_Marshal_Error   = 10003 // 資料Unmarshal錯誤
	Code_Data_Error           = 10004 // 資料錯誤 Unmarshal 成功但內容不正確
	Code_Data_Is_Exist        = 10005 // 資料已經存在
	Code_Data_Not_Exist       = 10006 // 資料不存在
	Code_Timeout              = 10007 // 超時
	Code_Token_Not_Exist      = 10008 // token不存在
	Code_Handler_Not_Exist    = 10009 // handler不存在

	// admin

	// api
	Code_Account_Length_Error      = 30001 // account 長度錯誤
	Code_Password_Length_Error     = 30002 // password 長度錯誤
	Code_NickName_Length_Error     = 30003 // nickname 長度錯誤
	Code_Account_Password_Error    = 30004 // 帳號或密碼錯誤
	Code_Club_Not_Exist            = 30005 // 俱樂部不存在
	Code_Club_User_Not_Permissions = 30006 // 權限不足
	Code_Club_User_Not_In_Club     = 30007 // 會員不在俱樂部
	Code_Club_User_Balacne_Less    = 30008 // 玩家金額不足
	Code_Club_User_Account_Exist   = 30009 // 玩家帳號已經存在

	// game socket server

	// table server
	Code_Table_Server_Req_Error = 50001 // table server 請求錯誤
	Code_Table_Not_Exist        = 50002 // 牌桌不存在
	Code_Table_LV_Insufficient  = 50003 // 牌桌等級不足
	Code_Table_Status_Error     = 50004 // 牌桌狀態錯誤
	Code_Game_Not_Type          = 50005 // 沒有此遊戲類型
	Code_Game_ID_Error          = 50006 // 遊戲ID錯誤
	Code_Game_Setting_Error     = 50007 // 遊戲設定驗證錯誤

	// game server
	Code_Game_Server_Create_Room_Error = 60001 // gs 創建房間失敗
	Code_GameWallet_Not_Exist          = 60002 // 錢包不存在
	Code_GameWallet_Balacne_Less       = 60003 // 錢包金額不足

	// db
	Code_DB_Find_Error        = 70001 // db get 錯誤
	Code_DB_Cursor_Error      = 70002 // db get 錯誤
	Code_DB_Count_Error       = 70002 // db count 錯誤
	Code_DB_Insert_Error      = 70003 // db insert 錯誤
	Code_DB_InsertMany_Error  = 70004 // db insert 錯誤
	Code_DB_Update_Error      = 70005 // db update 錯誤
	Code_DB_Delete_Error      = 70006 // db dle 錯誤
	Code_DB_Transaction_Error = 70007 // db transaction 錯誤
	Code_DB_Aggregate_Error   = 70007 // db transaction 錯誤
	Code_DB_BulkWrite_Error   = 70007 // db transaction 錯誤

	// Redis
	Code_Redis_Error = 90001 //

	// mq
	Code_Mq_Error     = 100001 // mq err
	Code_MQ_Pub_Error = 100002 // mq publish 錯誤
	Code_MQ_Sub_Error = 100003 // mq subject 錯誤
)
