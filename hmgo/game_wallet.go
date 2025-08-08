package hmgo

// import (
// 	"context"

// 	"github.com/chu16537/module_master/errorcode"
// 	"github.com/chu16537/module_master/proto/db"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// // 取得錢包
// func (h *Handler) GetGameWallet(ctx context.Context, filter *db.GameWalletOpt, opts *db.FindOpt) ([]*db.GameWallet, *errorcode.Error) {
// 	cur, err := h.find(ctx, db.ColName_GameWallet, filter.Filter_Mgo(), opts)
// 	if err != nil {
// 		return nil, errorcode.Server(err)
// 	}

// 	defer cur.Close(ctx)

// 	var results []*db.GameWallet

// 	if err := cur.All(ctx, &results); err != nil {
// 		return nil, errorcode.Server(err)
// 	}

// 	return results, nil
// }

// // 取得count GameWallet
// func (h *Handler) GetGameWalletCount(ctx context.Context, filter *db.GameWalletOpt) (int64, *errorcode.Error) {
// 	count, err := h.count(ctx, db.ColName_GameWallet, filter.Filter_Mgo())
// 	if err != nil {
// 		return 0, errorcode.Server(err)
// 	}
// 	return count, nil
// }

// // 創建遊戲錢包
// func (h *Handler) createGameWallet(ctx context.Context, data *db.GameWallet) *errorcode.Error {
// 	err := h.create(ctx, db.ColName_GameWallet, data)
// 	if err != nil {
// 		return errorcode.Server(err)
// 	}

// 	return nil
// }

// // 更新 game wallet
// func (h *Handler) UpdateGameWallet(ctx context.Context, filter *db.GameWalletOpt, data map[string]interface{}) (int64, *errorcode.Error) {
// 	c, err := h.update(ctx, db.ColName_GameWallet, filter.Filter_Mgo(), data)
// 	if err != nil {
// 		return c, errorcode.Server(err)
// 	}

// 	return c, nil
// }

// // 刪除遊戲錢包
// func (h *Handler) delGameWallet(ctx context.Context, filter *db.GameWalletOpt) *errorcode.Error {
// 	col := h.write.GetDB().Collection(db.ColName_GameWallet)

// 	_, err := col.DeleteMany(ctx, filter.Filter_Mgo())
// 	if err != nil {
// 		return errorcode.Server(err)
// 	}

// 	return nil
// }

// // 更新遊戲錢包金額
// func (h *Handler) updateGameWalletBalance(ctx context.Context, datas []*db.GameWallet) *errorcode.Error {
// 	col := h.write.GetDB().Collection(db.ColName_GameWallet)

// 	wms := make([]mongo.WriteModel, len(datas))
// 	for i, v := range datas {
// 		filter := bson.D{{Key: "user_id", Value: v.UserID}, {Key: "table_id", Value: v.TableID}}
// 		update := bson.D{{Key: "$set", Value: bson.D{
// 			{Key: "balance", Value: v.Balance},
// 			{Key: "total_bet", Value: v.TotalBet},
// 			{Key: "total_win", Value: v.TotalWin},
// 			{Key: "total_effective_bet", Value: v.TotalEffectiveBet}}}}
// 		wms[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
// 	}

// 	if _, err := col.BulkWrite(ctx, wms); err != nil {
// 		return errorcode.Server(err)
// 	}

// 	return nil
// }
