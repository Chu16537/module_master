package hmgo

// import (
// 	"context"
// 	"fmt"

// 	"github.com/Chu16537/module_master/errorcode"
// 	"github.com/Chu16537/module_master/proto/db"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// /*
// 創建俱樂部

// 測試:完成
// */
// func (h *Handler) CreateClub(ctx context.Context, account string, password string, nickName string, clubName string, status int, expireTime int64, invitationCode string, token string) (*db.Club, *db.ClubUserInfo, *errorcode.Error) {
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

// 	// 事務流程
// 	f := func(trans *Handler, sctx mongo.SessionContext) (interface{}, *errorcode.Error) {
// 		// 取得id
// 		id, err := trans.incr(sctx, db.ColName_Club)
// 		if err != nil {
// 			return nil, err
// 		}

// 		club.ID = uint64(id)

// 		// 創建俱樂部
// 		err = trans.createClub(sctx, club)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 取得id
// 		userID, err := trans.incr(sctx, db.ColName_Club_User_Info)
// 		if err != nil {
// 			return nil, err
// 		}

// 		clubUserInfo.UserID = uint64(userID)
// 		clubUserInfo.ClubID = club.ID

// 		// 創建俱樂部會員資料
// 		err = trans.createClubUserInfo(sctx, clubUserInfo)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return nil, nil
// 	}

// 	err := h.startTransaction(ctx, f)

// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return club, clubUserInfo, errorcode.Success()
// }

// /*
// 加入俱樂部

// 測試:完成
// */
// func (h *Handler) JoinClub(ctx context.Context, account string, password string, nickname string, clubID uint64, token string) (*db.ClubUserInfo, *errorcode.Error) {
// 	cui := &db.ClubUserInfo{
// 		Account:     account,
// 		Password:    password,
// 		NickName:    nickname,
// 		ClubID:      clubID,
// 		Token:       token,
// 		Permissions: db.Club_User_Permissions_Member,
// 	}

// 	// 事務流程
// 	f := func(trans *Handler, sctx mongo.SessionContext) (interface{}, *errorcode.Error) {
// 		cuiOpt := &db.ClubUserInfoOpt{
// 			ClubID: clubID,
// 			OR: &db.ClubUserInfoOR{
// 				Account:  account,
// 				NickName: nickname,
// 			},
// 		}

// 		acs, err := trans.GetClubUserInfo(sctx, cuiOpt, nil)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 有資料 不能創建
// 		if len(acs) > 0 {
// 			return nil, errorcode.New(errorcode.Club_User_Account_Exist, fmt.Errorf("account:%v or nickname:%v is exist", account, nickname))
// 		}

// 		// 取得id
// 		userID, err := trans.incr(sctx, db.ColName_Club_User_Info)
// 		if err != nil {
// 			return nil, err
// 		}

// 		cui.UserID = uint64(userID)

// 		err = trans.createClubUserInfo(sctx, cui)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return nil, nil
// 	}

// 	err := h.startTransaction(ctx, f)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return cui, errorcode.Success()
// }

// /*
// 檢查使用者權限是否可以操作

// 測試:完成
// */
// func (h *Handler) IsClubUserInfoPresident(ctx context.Context, opt *db.ClubUserInfoOpt) (*db.ClubUserInfo, *errorcode.Error) {
// 	cui, err := h.GetClubUserInfo(ctx, opt, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 代表此人沒有權限
// 	if len(cui) <= 0 {
// 		return nil, errorcode.NotClubPermissions(opt.UserID, opt.ClubID)
// 	}

// 	return cui[0], nil
// }

// /*
// 俱樂部轉帳

// 會員不能互轉

// 測試:完成
// */
// func (h *Handler) TransBalanceClub(ctx context.Context, presidentUserID uint64, tagUserID uint64, clubID uint64, amount uint64, orderID string, createTime int64) (uint64, *errorcode.Error) {
// 	var nowBalance uint64
// 	// 事務流程
// 	f := func(trans *Handler, sctx mongo.SessionContext) (interface{}, *errorcode.Error) {
// 		// presidentUserID 權限判斷
// 		presidentUserOpt := &db.ClubUserInfoOpt{
// 			UserID:      presidentUserID,
// 			ClubID:      clubID,
// 			Permissions: []int{db.Club_User_Permissions_President, db.Club_User_Permissions_Vice_President, db.Club_User_Permissions_Member},
// 		}
// 		pUser, err := trans.IsClubUserInfoPresident(sctx, presidentUserOpt)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// tagUserID 權限判斷
// 		tagUserOpt := &db.ClubUserInfoOpt{
// 			UserID:      tagUserID,
// 			ClubID:      clubID,
// 			Permissions: []int{db.Club_User_Permissions_President, db.Club_User_Permissions_Vice_President, db.Club_User_Permissions_Member},
// 		}
// 		tUser, err := trans.IsClubUserInfoPresident(sctx, tagUserOpt)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 會員不能互轉
// 		if pUser.Permissions == db.Club_User_Permissions_Member && tUser.Permissions == db.Club_User_Permissions_Member {
// 			return nil, errorcode.NotClubPermissions(presidentUserID, clubID)
// 		}

// 		// 更新金額資料
// 		update := make([]*db.UpdateBalanceInfo, 2)

// 		var userAmount uint64 = 0
// 		// 不是會長 要判斷金額
// 		if pUser.Permissions != db.Club_User_Permissions_President {
// 			if pUser.Balance < amount {
// 				err = errorcode.ClubUserBalanceLess(presidentUserID, pUser.Balance, amount)
// 				return nil, err
// 			}

// 			userAmount = amount
// 		}

// 		update[0] = &db.UpdateBalanceInfo{
// 			UserID:  presidentUserID,
// 			ClubID:  clubID,
// 			Balance: pUser.Balance - userAmount,
// 		}

// 		nowBalance = update[0].Balance

// 		update[1] = &db.UpdateBalanceInfo{
// 			UserID:  tagUserID,
// 			ClubID:  clubID,
// 			Balance: tUser.Balance + amount,
// 		}

// 		// 更新俱樂部錢包
// 		err = trans.updateClubUserInfoBalance(sctx, update)
// 		if err != nil {
// 			return nil, err
// 		}

// 		walletInfo := &db.WalletInfoClub{
// 			OrderID:              orderID,
// 			UserID:               presidentUserID,
// 			UserAmount:           int64(-userAmount),
// 			UserBeforeBalance:    pUser.Balance,
// 			TagUsertID:           tagUserID,
// 			TagUserAmount:        int64(amount),
// 			TagUserBeforeBalance: tUser.Balance,
// 			CreateTime:           createTime,
// 			ClubID:               clubID,
// 		}
// 		// 寫log
// 		logs := walletInfo.Log()
// 		err = trans.createWalletLog(sctx, logs)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return nil, nil
// 	}

// 	err := h.startTransaction(ctx, f)

// 	if err != nil {
// 		return 0, err
// 	}

// 	return nowBalance, errorcode.Success()
// }

// /*
// 創建錢包

// 玩家進入牌桌使用

// 測試:完成
// */
// func (h *Handler) CreateGameWallet(ctx context.Context, userID uint64, clubID uint64, tableID uint64, amount uint64, delTime int64, orderID string, createTime int64) (*db.GameWallet, *errorcode.Error) {
// 	gw := &db.GameWallet{
// 		UserID:  userID,
// 		ClubID:  clubID,
// 		TableID: tableID,
// 		DelTime: delTime,
// 	}

// 	// 事務流程
// 	f := func(trans *Handler, sctx mongo.SessionContext) (interface{}, *errorcode.Error) {
// 		// 查詢牌桌
// 		tOpt := &db.TableOpt{
// 			ID: tableID,
// 		}

// 		ts, err := trans.GetTable(sctx, tOpt, nil)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 牌桌不存在
// 		if len(ts) <= 0 {
// 			return nil, errorcode.TableNotExist(clubID, tableID)
// 		}

// 		// 取得俱樂部會員資料
// 		cuiOpt := &db.ClubUserInfoOpt{
// 			UserID:      userID,
// 			ClubID:      ts[0].ClubID,
// 			Permissions: []int{db.Club_User_Permissions_President, db.Club_User_Permissions_Vice_President, db.Club_User_Permissions_Member},
// 		}

// 		cui, err := trans.GetClubUserInfo(sctx, cuiOpt, nil)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if len(cui) <= 0 {
// 			return nil, errorcode.NotInClub(userID, ts[0].ClubID)
// 		}

// 		var updateClubAmount uint64 = 0
// 		// 不是會長 要判斷金額是否足夠
// 		if cui[0].Permissions != db.Club_User_Permissions_President {
// 			if cui[0].Balance < amount {
// 				return nil, errorcode.ClubUserBalanceLess(userID, cui[0].Balance, amount)
// 			}

// 			updateClubAmount = amount
// 		}

// 		// clubUserInfo 扣款
// 		updateClubUserInfoBalanceInfo := make([]*db.UpdateBalanceInfo, 1)
// 		updateClubUserInfoBalanceInfo[0] = &db.UpdateBalanceInfo{
// 			UserID:  userID,
// 			ClubID:  ts[0].ClubID,
// 			Balance: cui[0].Balance - updateClubAmount,
// 		}

// 		// 更新俱樂部金額
// 		err = trans.updateClubUserInfoBalance(sctx, updateClubUserInfoBalanceInfo)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 取得錢包
// 		gwOpt := &db.GameWalletOpt{
// 			UserID:  userID,
// 			TableID: tableID,
// 		}

// 		gws, err := trans.GetGameWallet(sctx, gwOpt, nil)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 沒有錢包 創建空錢包
// 		if len(gws) == 0 {
// 			gw.ClubID = ts[0].ClubID

// 			err = trans.createGameWallet(sctx, gw)
// 			if err != nil {
// 				return nil, err
// 			}
// 		} else {
// 			gw = gws[0]
// 		}

// 		// 更新錢包金額
// 		updateGameWalletBalanceInfo := make([]*db.GameWallet, 1)
// 		updateGameWalletBalanceInfo[0] = &db.GameWallet{
// 			UserID:  userID,
// 			TableID: tableID,
// 			Balance: gw.Balance + amount,
// 		}

// 		err = trans.updateGameWalletBalance(sctx, updateGameWalletBalanceInfo)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 創建log
// 		walletInfo := &db.WalletInfoClubToGame{
// 			OrderID:           orderID,
// 			UserID:            userID,
// 			UpdateGameAmount:  int64(amount),
// 			GameBeforeBalance: gw.Balance,
// 			UpdateClubAmount:  -int64(updateClubAmount),
// 			ClubBeforeBalance: cui[0].Balance,
// 			CreateTime:        createTime,
// 			ClubID:            ts[0].ClubID,
// 			TableID:           tableID,
// 		}

// 		err = trans.createWalletLog(sctx, walletInfo.Log())
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 更新金額
// 		gw.Balance = updateGameWalletBalanceInfo[0].Balance

// 		return nil, nil
// 	}

// 	err := h.startTransaction(ctx, f)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return gw, errorcode.Success()
// }

// /*
// 玩家在遊戲中更新錢包金額

// 提款 amount < 0 從 遊戲 匯款到 俱樂部

// 匯入 amount > 0 從 俱樂部 匯款到 遊戲

// 測試:完成
// */
// func (h *Handler) UpdateGameWalletBalance(ctx context.Context, userID uint64, clubID uint64, tableID uint64, amount int64, orderID string, createTime int64) (uint64, uint64, *errorcode.Error) {
// 	var (
// 		newClubBalance uint64 = 0
// 		newGameBalance uint64 = 0
// 	)

// 	// 事務流程
// 	f := func(trans *Handler, sctx mongo.SessionContext) (interface{}, *errorcode.Error) {
// 		// 查詢牌桌
// 		tOpt := &db.TableOpt{
// 			ID:     tableID,
// 			ClubID: clubID,
// 		}

// 		ts, err := trans.GetTable(sctx, tOpt, nil)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 牌桌不存在
// 		if len(ts) <= 0 {
// 			return nil, errorcode.TableNotExist(clubID, tableID)
// 		}

// 		// 取得俱樂部會員資料
// 		cuiOpt := &db.ClubUserInfoOpt{
// 			UserID:      userID,
// 			ClubID:      ts[0].ClubID,
// 			Permissions: []int{db.Club_User_Permissions_President, db.Club_User_Permissions_Vice_President, db.Club_User_Permissions_Member},
// 		}

// 		cui, err := trans.IsClubUserInfoPresident(sctx, cuiOpt)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 取得錢包
// 		gwOpt := &db.GameWalletOpt{
// 			UserID:  userID,
// 			ClubID:  ts[0].ClubID,
// 			TableID: tableID,
// 		}

// 		gws, err := trans.GetGameWallet(sctx, gwOpt, nil)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if len(gws) <= 0 {
// 			return nil, errorcode.GameWalletNotExist(userID, ts[0].ID)
// 		}

// 		var (
// 			updateClubAmount     int64 = -amount
// 			updateGameClubAmount int64 = amount
// 			logs                 []*db.WalletLog
// 		)

// 		if amount < 0 {
// 			// 從 遊戲 匯款到 俱樂部

// 			// 遊戲錢包不夠
// 			if gws[0].Balance < uint64(-amount) {
// 				return nil, errorcode.GameWalleBalanceLess(userID, gws[0].TableID, gws[0].Balance, uint64(-amount))
// 			}

// 			// 創建log
// 			walletInfo := &db.WalletInfoGameToClub{
// 				OrderID:           orderID,
// 				UserID:            userID,
// 				UpdateGameAmount:  updateGameClubAmount,
// 				GameBeforeBalance: gws[0].Balance,
// 				UpdateClubAmount:  updateClubAmount,
// 				ClubBeforeBalance: cui.Balance,
// 				CreateTime:        createTime,
// 				ClubID:            ts[0].ClubID,
// 				TableID:           tableID,
// 			}

// 			logs = walletInfo.Log()
// 		} else if amount > 0 {
// 			// 從 俱樂部 匯款到 遊戲

// 			// 不是會長需要判斷金額
// 			if cui.Permissions != db.Club_User_Permissions_President {
// 				// 遊俱樂部錢包不夠
// 				if cui.Balance < uint64(amount) {
// 					return nil, errorcode.ClubUserBalanceLess(userID, cui.Balance, uint64(amount))
// 				}
// 			} else if cui.Permissions == db.Club_User_Permissions_President {
// 				updateClubAmount = 0
// 			}

// 			walletInfo := &db.WalletInfoClubToGame{
// 				OrderID:           orderID,
// 				UserID:            userID,
// 				UpdateGameAmount:  updateGameClubAmount,
// 				GameBeforeBalance: gws[0].Balance,
// 				UpdateClubAmount:  updateClubAmount,
// 				ClubBeforeBalance: cui.Balance,
// 				CreateTime:        createTime,
// 				ClubID:            ts[0].ClubID,
// 				TableID:           tableID,
// 			}

// 			logs = walletInfo.Log()
// 		}

// 		newClubBalance = uint64(int64(cui.Balance) + updateClubAmount)
// 		newGameBalance = uint64(int64(gws[0].Balance) + updateGameClubAmount)

// 		// clubUserInfo 金額更新
// 		updateClubUserInfoBalanceInfo := make([]*db.UpdateBalanceInfo, 1)
// 		updateClubUserInfoBalanceInfo[0] = &db.UpdateBalanceInfo{
// 			UserID:  userID,
// 			ClubID:  ts[0].ClubID,
// 			Balance: newClubBalance,
// 		}

// 		err = trans.updateClubUserInfoBalance(sctx, updateClubUserInfoBalanceInfo)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// game wallet 金額更新
// 		updateGameWalletBalanceInfo := make([]*db.GameWallet, 1)
// 		updateGameWalletBalanceInfo[0] = &db.GameWallet{
// 			UserID:  userID,
// 			TableID: tableID,
// 			Balance: newGameBalance,
// 		}

// 		err = trans.updateGameWalletBalance(sctx, updateGameWalletBalanceInfo)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 創建log
// 		err = trans.createWalletLog(sctx, logs)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return nil, nil
// 	}

// 	err := h.startTransaction(ctx, f)

// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	return newClubBalance, newGameBalance, errorcode.Success()
// }

// /*
// 刪除錢包

// 測試:完成
// */
// func (h *Handler) DeleteGameWallet(ctx context.Context, gwOpt *db.GameWalletOpt, orderID string, createTime int64) *errorcode.Error {
// 	// 事務流程
// 	f := func(trans *Handler, sctx mongo.SessionContext) (interface{}, *errorcode.Error) {
// 		// 取得錢包
// 		gws, err := trans.GetGameWallet(sctx, gwOpt, nil)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 沒有錢包 不做事
// 		if len(gws) <= 0 {
// 			return nil, nil
// 		}

// 		// 轉為map 方便查找
// 		gwMap := map[uint64]*db.GameWallet{}

// 		for _, v := range gws {
// 			gwMap[v.UserID] = v
// 		}

// 		userIDs := make([]uint64, len(gws))
// 		for i, v := range gws {
// 			userIDs[i] = v.UserID
// 		}

// 		// 取得俱樂部會員資料
// 		cuiOpt := &db.ClubUserInfoOpt{
// 			UserIDs: userIDs,
// 		}

// 		cuis, err := trans.GetClubUserInfo(sctx, cuiOpt, nil)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 轉為map 方便查找
// 		cuiMap := map[uint64]*db.ClubUserInfo{}
// 		for _, v := range cuis {
// 			cuiMap[v.UserID] = v
// 		}

// 		updateClubUserInfoBalanceInfo := make([]*db.UpdateBalanceInfo, len(userIDs))
// 		cIdx := 0
// 		logs := make([]*db.WalletLog, len(userIDs)*2)
// 		lIdx := 0

// 		for uID, gameWallet := range gwMap {
// 			userInfo, ok := cuiMap[uID]
// 			// 玩家不存在 跳過
// 			if !ok {
// 				continue
// 			}

// 			// clubUserInfo 金額更新資料
// 			updateClubUserInfoBalanceInfo[cIdx] = &db.UpdateBalanceInfo{
// 				UserID:  uID,
// 				ClubID:  userInfo.ClubID,
// 				Balance: userInfo.Balance + gameWallet.Balance,
// 			}

// 			cIdx++

// 			// 創建log
// 			walletInfo := &db.WalletInfoGameToClub{
// 				OrderID:           orderID,
// 				UserID:            uID,
// 				UpdateGameAmount:  int64(-gameWallet.Balance),
// 				GameBeforeBalance: gameWallet.Balance,
// 				UpdateClubAmount:  int64(gameWallet.Balance),
// 				ClubBeforeBalance: userInfo.Balance,
// 				CreateTime:        createTime,
// 				ClubID:            userInfo.ClubID,
// 				TableID:           gameWallet.TableID,
// 			}

// 			ls := walletInfo.Log()
// 			logs[lIdx] = ls[0]
// 			lIdx++
// 			logs[lIdx] = ls[1]
// 			lIdx++
// 		}

// 		// clubUserInfo 金額更新
// 		updateClubUserInfoBalanceInfo = updateClubUserInfoBalanceInfo[:cIdx]
// 		err = trans.updateClubUserInfoBalance(sctx, updateClubUserInfoBalanceInfo)
// 		if err != nil {
// 			return nil, err
// 		}

// 		err = trans.delGameWallet(sctx, gwOpt)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 創建log
// 		logs = logs[:lIdx]
// 		err = trans.createWalletLog(sctx, logs)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return nil, nil
// 	}

// 	err := h.startTransaction(ctx, f)

// 	if err != nil {
// 		return err
// 	}

// 	return errorcode.Success()
// }

// /*
// 創建注單 並 更新遊戲錢包

// 1. 創建GameRecor

// 2. 更新錢包

// 3. 錢包log 下注 跟 派獎 兩個

// 4. 創建 UserRecord

// 測試:完成
// */
// func (h *Handler) CreateGameRecord(ctx context.Context, gr *db.GameRecord, urs []*db.UserRecord, info map[uint64]*db.UserRecordInfodMultiple) *errorcode.Error {
// 	// 事務流程
// 	f := func(trans *Handler, sctx mongo.SessionContext) (interface{}, *errorcode.Error) {
// 		// 創建gr
// 		err := trans.createGameRecord(sctx, gr)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 代表不用更新 UserRecord 跟 錢包
// 		if len(urs) == 0 || len(info) == 0 {
// 			return nil, nil
// 		}

// 		// 取得錢包
// 		gameWalletOpt := &db.GameWalletOpt{
// 			TableID: gr.TableID,
// 		}

// 		gws, err := trans.GetGameWallet(sctx, gameWalletOpt, nil)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if len(gws) == 0 {
// 			return nil, nil
// 		}

// 		infoIdx := 0
// 		updateBalanceInfos := make([]*db.GameWallet, len(gws))
// 		logIdx := 0
// 		logs := make([]*db.WalletLog, len(gws)*2)

// 		for _, v := range gws {
// 			// 不存在
// 			if _, ok := info[v.UserID]; !ok {
// 				continue
// 			}

// 			// 當前金額下注錯誤
// 			// 理論上不會在這邊報錯 因為 下注時就會判斷是否正確
// 			if v.Balance-info[v.UserID].TotalBet <= 0 {
// 				continue
// 			}

// 			updateBalanceInfos[infoIdx] = &db.GameWallet{
// 				UserID:            v.UserID,
// 				TableID:           gr.TableID,
// 				Balance:           v.Balance - info[v.UserID].TotalBet + info[v.UserID].TotalWin,
// 				TotalBet:          v.TotalBet + info[v.UserID].TotalBet,
// 				TotalWin:          v.TotalWin + info[v.UserID].TotalWin,
// 				TotalEffectiveBet: v.TotalEffectiveBet + info[v.UserID].TotalEffectiveBet,
// 			}

// 			infoIdx++

// 			// 創建log
// 			walletInfo := &db.WalletInfoGameRecode{
// 				OrderID:       gr.GameRecordID,
// 				UserID:        v.UserID,
// 				Bet:           info[v.UserID].TotalBet,
// 				Win:           info[v.UserID].TotalWin,
// 				BeforeBalance: v.Balance,
// 				CreateTime:    gr.CreateTime,
// 				ClubID:        gr.ClubID,
// 				TableID:       gr.TableID,
// 			}

// 			ws := walletInfo.Log()
// 			logs[logIdx] = ws[0]
// 			logIdx++
// 			logs[logIdx] = ws[1]
// 			logIdx++
// 		}

// 		// 更新錢包
// 		updateBalanceInfos = updateBalanceInfos[:infoIdx]
// 		err = trans.updateGameWalletBalance(sctx, updateBalanceInfos)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// 寫log
// 		logs = logs[:logIdx]
// 		err = trans.createWalletLog(sctx, logs)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if len(urs) != 0 {
// 			err = trans.CreateUserRecord(sctx, urs)
// 			if err != nil {
// 				return nil, err
// 			}
// 		}

// 		return nil, nil
// 	}

// 	err := h.startTransaction(ctx, f)

// 	if err != nil {
// 		return err
// 	}

// 	return errorcode.Success()
// }
