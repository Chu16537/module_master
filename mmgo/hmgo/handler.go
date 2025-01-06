package hmgo

// // 創建俱樂部
// // 創建會長資料
// func CreateClub(mgo *mmgo.Handler, ctx context.Context, account string, password string, nickName string, clubName string, status int, expireTime int64, invitationCode string, token string) (*db.Club, *db.ClubUserInfo, *errorcode.Error) {
// 	club := &db.Club{
// 		Name:           clubName,
// 		PresidentName:  nickName,
// 		Status:         status,
// 		ExpireTime:     expireTime,
// 		InvitationCode: invitationCode,
// 	}

// 	clubUserInfo := &db.ClubUserInfo{
// 		Account:     account,
// 		Password:    password,
// 		NickName:    nickName,
// 		Permissions: db.Club_User_Permissions_President,
// 		Token:       token,
// 	}

// 	// 指定好讀寫db
// 	h := mgo.GetWW()

// 	// 取得db
// 	database := h.GetDBWrite()

// 	//  執行 Transaction
// 	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
// 		// 創建俱樂部
// 		err := create(database, sessCtx, db.ColName_Club, club)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 創建俱樂部會長
// 		err = create(database, sessCtx, db.ColName_Club_User_Info, clubUserInfo)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return nil, nil
// 	}

// 	// 執行 Transaction
// 	err := h.RunTransaction(ctx, callback)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return club, clubUserInfo, errorcode.Success()
// }

// // 加入俱樂部
// func JoinClub(mgo *mmgo.Handler, ctx context.Context, account string, password string, nickname string, clubID uint64, token string) (*db.ClubUserInfo, *errorcode.Error) {
// 	clubUserInfo := &db.ClubUserInfo{
// 		Account:     account,
// 		Password:    password,
// 		NickName:    nickname,
// 		ClubID:      clubID,
// 		Token:       token,
// 		Permissions: db.Club_User_Permissions_Member,
// 	}

// 	// 指定好讀寫db
// 	h := mgo.GetWW()

// 	// 取得db
// 	database := h.GetDBWrite()

// 	// 執行 Transaction
// 	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
// 		// 查詢是否有相同帳號
// 		cuiOpt := &db.ClubUserInfoOpt{
// 			ClubID:  clubID,
// 			Account: account,
// 		}

// 		cuis, err := GetClubUserInfo(database, sessCtx, cuiOpt, nil)
// 		// 有資料 不能創建
// 		if len(cuis) > 0 && err == nil {
// 			err = fmt.Errorf("account:%v or nickname:%v is exist", account, nickname)
// 		}
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 創建俱樂部會員
// 		err = create(database, sessCtx, db.ColName_Club_User_Info, clubUserInfo)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return nil, nil
// 	}

// 	// 執行 Transaction
// 	err := h.RunTransaction(ctx, callback)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return clubUserInfo, errorcode.Success()
// }
