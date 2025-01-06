所有func 都要回傳 *errorcode.Error

即使沒有錯誤也要回傳 errorcode.Success()

以防止panic 發生