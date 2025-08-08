package hmgo

// import (
// 	"context"

// 	"github.com/chu16537/module_master/errorcode"
// 	"github.com/chu16537/module_master/proto/db"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// 取得 ClubUserInfo
// func (h *Handler) GetClubUserInfo(ctx context.Context, filter *db.ClubUserInfoOpt, opts *db.FindOpt) ([]*db.ClubUserInfo, *errorcode.Error) {
// 	cur, err := h.find(ctx, db.ColName_Club_User_Info, filter.Filter_Mgo(), opts)
// 	if err != nil {
// 		return nil, errorcode.Server(err)
// 	}

// 	defer cur.Close(ctx)

// 	var datas []*db.ClubUserInfo

// 	if err := cur.All(ctx, &datas); err != nil {
// 		return nil, errorcode.Server(err)
// 	}

// 	return datas, nil
// }

// // 取得count ClubUserInfo
// func (h *Handler) GetClubUserInfoCount(ctx context.Context, filter *db.ClubUserInfoOpt) (int64, *errorcode.Error) {
// 	count, err := h.count(ctx, db.ColName_Club_User_Info, filter.Filter_Mgo())
// 	if err != nil {
// 		return 0, errorcode.Server(err)
// 	}
// 	return count, nil
// }

// // 創建 ClubUserInfo
// func (h *Handler) createClubUserInfo(ctx context.Context, data *db.ClubUserInfo) *errorcode.Error {
// 	err := h.create(ctx, db.ColName_Club_User_Info, data)
// 	if err != nil {
// 		return errorcode.Server(err)
// 	}

// 	return nil
// }

// // 更新 ClubUserInfo
// func (h *Handler) UpdateClubUserInfo(ctx context.Context, filter *db.ClubUserInfoOpt, data map[string]interface{}) *errorcode.Error {
// 	_, err := h.update(ctx, db.ColName_Club_User_Info, filter.Filter_Mgo(), data)
// 	if err != nil {
// 		return errorcode.Server(err)
// 	}

// 	return errorcode.Success()
// }

// // 更新金額
// func (h *Handler) updateClubUserInfoBalance(ctx context.Context, datas []*db.UpdateBalanceInfo) *errorcode.Error {
// 	col := h.write.GetDB().Collection(db.ColName_Club_User_Info)

// 	wms := make([]mongo.WriteModel, len(datas))
// 	for i, v := range datas {
// 		filter := bson.D{{Key: "user_id", Value: v.UserID}, {Key: "club_id", Value: v.ClubID}}
// 		update := bson.D{{Key: "$set", Value: bson.D{{Key: "balance", Value: v.Balance}}}}
// 		wms[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
// 	}

// 	_, err := col.BulkWrite(ctx, wms)
// 	if err != nil {
// 		return errorcode.Server(err)
// 	}

// 	return nil
// }
