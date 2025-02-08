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

	// db
	Code_DB_Find_Error        = 20001 // db get 錯誤
	Code_DB_Cursor_Error      = 20002 // db get 錯誤
	Code_DB_Count_Error       = 20002 // db count 錯誤
	Code_DB_Insert_Error      = 20003 // db insert 錯誤
	Code_DB_InsertMany_Error  = 20004 // db insert 錯誤
	Code_DB_Update_Error      = 20005 // db update 錯誤
	Code_DB_Delete_Error      = 20006 // db dle 錯誤
	Code_DB_Transaction_Error = 20007 // db transaction 錯誤
	Code_DB_Aggregate_Error   = 20007 // db transaction 錯誤
	Code_DB_BulkWrite_Error   = 20007 // db transaction 錯誤

	// Redis
	Code_Redis_Error = 30001 //

	// mq
	Code_Mq_Error     = 400001 // mq err
	Code_MQ_Pub_Error = 400002 // mq publish 錯誤
	Code_MQ_Sub_Error = 400003 // mq subject 錯誤

	// 以下為各服務的錯誤
	// admin

	// api
	Code_Platform_Not_Exist = 100001 // Platform 不存在
	Code_Decrypt_Error      = 100002 // Decrypt 錯誤

	// game socket server

	// table server
	Code_Table_Server_Req_Error = 110001 // table server 請求錯誤
	Code_Table_Not_Exist        = 110002 // 牌桌不存在
	Code_Table_LV_Insufficient  = 110003 // 牌桌等級不足
	Code_Table_Status_Error     = 110004 // 牌桌狀態錯誤
	Code_Game_Not_Type          = 110005 // 沒有此遊戲類型
	Code_Game_ID_Error          = 110006 // 遊戲ID錯誤
	Code_Game_Setting_Error     = 110007 // 遊戲設定驗證錯誤

	// game server
	Code_Game_Server_Create_Room_Error = 120001 // gs 創建房間失敗
	Code_GameWallet_Not_Exist          = 120002 // 錢包不存在
	Code_GameWallet_Balacne_Less       = 120003 // 錢包金額不足

)
