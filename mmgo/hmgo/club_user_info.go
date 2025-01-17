package hmgo

import (
	"context"
	"fmt"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mmgo"
	"github.com/Chu16537/module_master/proto/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 取得俱樂部會員資料
func GetClubUserInfo(mh *mmgo.Handler, ctx context.Context, cuiOpt *db.ClubUserInfoOpt, findOpt *db.FindOpt) ([]*db.ClubUserInfo, *errorcode.Error) {
	datas, err := mmgo.Find[*db.ClubUserInfo](mh, ctx, db.ColName_Club_User_Info, cuiOpt, findOpt)

	return datas, err
}

// 加入俱樂部
func CreateClubUserInfo(mh *mmgo.Handler, ctx context.Context, account string, password string, nickname string, clubID uint64, token string) (*db.ClubUserInfo, *errorcode.Error) {
	id, err := mmgo.Incr(mh, ctx, db.ColName_Club_User_Info)
	if err != nil {
		return nil, err
	}

	cui := &db.ClubUserInfo{
		UserID:      uint64(id),
		Account:     account,
		Password:    password,
		NickName:    nickname,
		ClubID:      clubID,
		Token:       token,
		Permissions: db.Club_User_Permissions_Member,
	}

	err = mmgo.Insert(mh, ctx, db.ColName_Club_User_Info, cui)

	return cui, err
}

// 更改密碼
func UpdateClubUserInfo(mh *mmgo.Handler, ctx context.Context, cuiOpt *db.ClubUserInfoOpt, data map[string]interface{}) *errorcode.Error {
	updateCount, err := mmgo.Update(mh, ctx, db.ColName_Club_User_Info, cuiOpt, data)

	if updateCount == 0 && !!err.IsSuccess() {
		err = errorcode.New(errorcode.Code_DB_Update_Error, fmt.Errorf("account:%v not exist", cuiOpt.Account))
	}

	return err
}

// 俱樂部會員轉帳
func TransBalanceClub(mh *mmgo.Handler, ctx context.Context, promoterUserID uint64, tagUserID uint64, clubID uint64, amount uint64, orderID string, createTime int64) (uint64, *errorcode.Error) {

	trans := mh.GetWW()
	var (
		nowBalance uint64 = 0
		ec         *errorcode.Error
	)

	fn := func(sctx mongo.SessionContext) (interface{}, error) {

		// promoterUserID 權限判斷
		presidentUserOpt := &db.ClubUserInfoOpt{
			UserID:      promoterUserID,
			ClubID:      clubID,
			Permissions: []int{db.Club_User_Permissions_President, db.Club_User_Permissions_Vice_President, db.Club_User_Permissions_Member},
		}

		pUser, err := mmgo.FindOne[*db.ClubUserInfo](trans, ctx, db.ColName_Club_User_Info, presidentUserOpt)
		if pUser == nil && !err.IsSuccess() {
			err = errorcode.New(errorcode.Code_DB_Find_Error, fmt.Errorf("promoterUserID:%v not exist", promoterUserID))
		}
		if !err.IsSuccess() {
			ec = err
			return nil, err.GetErr()
		}

		// tagUserID 權限判斷
		tagUserOpt := &db.ClubUserInfoOpt{
			UserID:      tagUserID,
			ClubID:      clubID,
			Permissions: []int{db.Club_User_Permissions_President, db.Club_User_Permissions_Vice_President, db.Club_User_Permissions_Member},
		}
		tUser, err := mmgo.FindOne[*db.ClubUserInfo](trans, ctx, db.ColName_Club_User_Info, tagUserOpt)
		if tUser == nil && !err.IsSuccess() {
			err = errorcode.New(errorcode.Code_DB_Find_Error, fmt.Errorf("tagUserID:%v not exist", tagUserID))
		}
		if !err.IsSuccess() {
			ec = err
			return nil, err.GetErr()
		}

		// 會員不能互轉
		if pUser.Permissions == db.Club_User_Permissions_Member && tUser.Permissions == db.Club_User_Permissions_Member {
			err = errorcode.New(errorcode.Code_Club_User_Not_Permissions, fmt.Errorf("promoterUser and tagUser permissions is member"))
			ec = err
			return nil, err.GetErr()
		}

		// 更新金額資料
		update := make([]*db.UpdateBalanceInfo, 2)

		var userAmount uint64 = 0
		// 不是會長 要判斷金額
		if pUser.Permissions != db.Club_User_Permissions_President {
			if pUser.Balance < amount {
				err = errorcode.New(errorcode.Code_Club_User_Balacne_Less, fmt.Errorf("presidentUserID:%v balance:%v less than amount:%v", promoterUserID, pUser.Balance, amount))
				ec = err
				return nil, err.GetErr()
			}

			userAmount = amount
		}

		nowBalance = pUser.Balance - userAmount

		update[0] = &db.UpdateBalanceInfo{
			UserID:  promoterUserID,
			ClubID:  clubID,
			Balance: nowBalance,
		}

		update[1] = &db.UpdateBalanceInfo{
			UserID:  tagUserID,
			ClubID:  clubID,
			Balance: tUser.Balance + amount,
		}

		wms := createUpdateClubUserInfoBalanceData(update)

		// 更新金額
		_, err = mmgo.BulkWrite(trans, sctx, db.ColName_Club_User_Info, wms)
		if !err.IsSuccess() {

			ec = err
			return nil, err.GetErr()
		}

		// 新增log
		walletInfo := &db.WalletInfoClub{
			OrderID:              orderID,
			UserID:               promoterUserID,
			UserAmount:           int64(-userAmount),
			UserBeforeBalance:    pUser.Balance,
			TagUsertID:           tagUserID,
			TagUserAmount:        int64(amount),
			TagUserBeforeBalance: tUser.Balance,
			CreateTime:           createTime,
			ClubID:               clubID,
		}
		fmt.Println("walletInfo", walletInfo)
		// 寫log
		logs := walletInfo.Log()
		fmt.Println("logs", logs)
		err = mmgo.InsertMany(trans, sctx, db.ColName_WalletLog, logs)
		if err != nil {
			ec = err
			return nil, err.GetErr()
		}
		return nil, nil
	}

	err := mmgo.Transaction(trans, ctx, fn)
	if ec != nil {
		return 0, ec
	}

	if err != nil {
		return 0, err
	}

	return nowBalance, errorcode.Success()
}

// 創建更新金額資料
func createUpdateClubUserInfoBalanceData(datas []*db.UpdateBalanceInfo) []mongo.WriteModel {
	wms := make([]mongo.WriteModel, len(datas))

	for i, v := range datas {
		filter := bson.D{{Key: "user_id", Value: v.UserID}, {Key: "club_id", Value: v.ClubID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "balance", Value: v.Balance}}}}
		wms[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
	}
	return wms
}
