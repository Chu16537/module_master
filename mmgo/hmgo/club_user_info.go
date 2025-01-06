package hmgo

import (
	"context"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mmgo"
	"github.com/Chu16537/module_master/proto/db"
)

// 取得 ClubUserInfo
func GetClubUserInfo(h *mmgo.Handler, ctx context.Context, filter *db.ClubUserInfoOpt, opts *db.FindOpt) ([]*db.ClubUserInfo, *errorcode.Error) {
	datas, err := mmgo.Find[*db.ClubUserInfo](h, ctx, db.ColName_Club_User_Info, filter.ToMgo(), opts)
	if err != nil {
		return nil, err
	}
	return datas, errorcode.Success()
}

// // // 取得count ClubUserInfo
// // func GetCountClubUserInfo(db *mongo.Database, ctx context.Context, filter *db.ClubUserInfoOpt) (int64, error) {
// // 	count, err := count(db, ctx, db.ColName_Club_User_Info, filter.Filter_Mgo())
// // 	if err != nil {
// // 		return 0, err
// // 	}
// // 	return count, nil
// // }

// 創建 ClubUserInfo
func CreateClubUserInfo(h *mmgo.Handler, ctx context.Context, data *db.ClubUserInfo) *errorcode.Error {
	err := mmgo.Create(h, ctx, db.ColName_Club_User_Info, data)
	if err != nil {
		return err
	}

	return errorcode.Success()
}

// // // // 更新 ClubUserInfo
// // // func (h *Handler) UpdateClubUserInfo(ctx context.Context, filter *db.ClubUserInfoOpt, data map[string]interface{}) *errorcode.Error {
// // // 	_, err := h.update(ctx, db.ColName_Club_User_Info, filter.Filter_Mgo(), data)
// // // 	if err != nil {
// // // 		return errorcode.Server(err)
// // // 	}

// // // 	return errorcode.Success()
// // // }

// // // // 更新金額
// // // func (h *Handler) updateClubUserInfoBalance(ctx context.Context, datas []*db.UpdateBalanceInfo) *errorcode.Error {
// // // 	col := h.write.GetDB().Collection(db.ColName_Club_User_Info)

// // // 	wms := make([]mongo.WriteModel, len(datas))
// // // 	for i, v := range datas {
// // // 		filter := bson.D{{Key: "user_id", Value: v.UserID}, {Key: "club_id", Value: v.ClubID}}
// // // 		update := bson.D{{Key: "$set", Value: bson.D{{Key: "balance", Value: v.Balance}}}}
// // // 		wms[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
// // // 	}

// // // 	_, err := col.BulkWrite(ctx, wms)
// // // 	if err != nil {
// // // 		return errorcode.Server(err)
// // // 	}

// // // 	return nil
// // // }
