package hmgo

// import (
// 	"context"

// 	"github.com/Chu16537/module_master/errorcode"
// 	"github.com/Chu16537/module_master/proto/db"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// // 取得俱樂部
// func (h *Handler) GetClub(ctx context.Context, filter *db.ClubOpt, opts *db.FindOpt) ([]*db.Club, *errorcode.Error) {
// 	cur, err := h.find(ctx, db.ColName_Club, filter.Filter_Mgo(), opts)
// 	if err != nil {
// 		return nil, errorcode.Server(err)
// 	}

// 	defer cur.Close(ctx)

// 	var results []*db.Club

// 	if err := cur.All(ctx, &results); err != nil {
// 		return nil, errorcode.Server(err)
// 	}

// 	return results, nil
// }

// // 取得count Club
// func (h *Handler) GetClubCount(ctx context.Context, filter *db.ClubOpt) (int64, *errorcode.Error) {
// 	count, err := h.count(ctx, db.ColName_Club, filter.Filter_Mgo())
// 	if err != nil {
// 		return 0, errorcode.Server(err)
// 	}
// 	return count, nil
// }

// // 創建俱樂部
// func (h *Handler) createClub(ctx context.Context, data *db.Club) *errorcode.Error {
// 	err := h.create(ctx, db.ColName_Club, data)
// 	if err != nil {
// 		return errorcode.Server(err)
// 	}

// 	return nil
// }

// // 更新俱樂部 / 回傳 更改的資料數量
// func (h *Handler) UpdateClub(ctx context.Context, filter *db.ClubOpt, data map[string]interface{}) (int64, *errorcode.Error) {
// 	c, err := h.update(ctx, db.ColName_Club, filter.Filter_Mgo(), data)
// 	if err != nil {
// 		return c, errorcode.Server(err)
// 	}

// 	return c, nil
// }

// // 使用 userid 取得俱樂部資料
// func (h *Handler) GetClubByUserId(ctx context.Context, userID uint64, clubID uint64) (*db.Club, *db.ClubUserInfo, *errorcode.Error) {
// 	col := h.read.GetDB().Collection(db.ColName_Club_User_Info)

// 	lookupStage := bson.D{
// 		{Key: "$lookup", Value: bson.D{
// 			{Key: "from", Value: db.ColName_Club},
// 			{Key: "localField", Value: "club_id"},
// 			{Key: "foreignField", Value: "id"},
// 			{Key: "as", Value: "club"},
// 		}},
// 	}

// 	matchStage := bson.D{
// 		{Key: "$match", Value: bson.D{
// 			{Key: "user_id", Value: userID},
// 			{Key: "club_id", Value: clubID},
// 		}},
// 	}

// 	unwindStage := bson.D{
// 		{Key: "$unwind", Value: "$club"},
// 	}

// 	projectStage := bson.D{
// 		{Key: "$project", Value: bson.D{
// 			{Key: "club_user_info", Value: "$$ROOT"},
// 			{Key: "club", Value: "$club"},
// 		}},
// 	}

// 	pipeline := mongo.Pipeline{lookupStage, matchStage, unwindStage, projectStage}

// 	cursor, err := col.Aggregate(ctx, pipeline)
// 	if err != nil {
// 		return nil, nil, errorcode.Server(err)
// 	}
// 	defer cursor.Close(ctx)

// 	var results []db.ClubByUserId

// 	if err := cursor.All(ctx, &results); err != nil {
// 		return nil, nil, errorcode.Server(err)
// 	}

// 	if len(results) == 0 {
// 		return nil, nil, errorcode.NotClubPermissions(userID, clubID)
// 	}

// 	return &results[0].Club, &results[0].ClubUserInfo, nil
// }
