package hmgo

// import (
// 	"context"

// 	"github.com/Chu16537/module_master/errorcode"
// 	"github.com/Chu16537/module_master/mtime"
// 	"github.com/Chu16537/module_master/proto/db"
// 	"go.mongodb.org/mongo-driver/bson"
// )

// // 取得 UserRecord
// func (h *Handler) GetUserRecord(ctx context.Context, filter *db.UserRecordOpt, opts *db.FindOpt) ([]*db.UserRecord, int64, *errorcode.Error) {
// 	dates := mtime.GetDateRange(filter.StartTimeUnix, filter.EndTimeUnix)

// 	pipeline := make([]bson.M, len(dates)+2)

// 	for i := range dates {
// 		pipeline[i] = bson.M{
// 			"$unionWith": bson.M{
// 				"coll": db.ColName_UserRecord + dates[i],
// 			},
// 		}
// 	}

// 	pipeline[len(pipeline)-2] = bson.M{
// 		"$match": filter.Filter_Mgo(),
// 	}

// 	pipeline[len(pipeline)-1] = bson.M{"$count": "count"}

// 	// 执行聚合查询以获取总数
// 	totalCountCursor, err := h.read.GetDB().Collection(db.ColName_GameRecord).Aggregate(ctx, pipeline)
// 	if err != nil {
// 		return nil, 0, errorcode.Server(err)
// 	}
// 	defer totalCountCursor.Close(ctx)

// 	totalCount := &db.TotalCount{}
// 	if totalCountCursor.Next(ctx) {
// 		if err := totalCountCursor.Decode(totalCount); err != nil {
// 			return nil, 0, errorcode.Server(err)
// 		}
// 	}

// 	pipeline = pipeline[:len(pipeline)-1]

// 	if opts != nil {
// 		opts.ToMgo()
// 		// 加入 $skip 跳過前面的紀錄
// 		pipeline = append(pipeline, bson.M{"$skip": opts.Start})
// 		// 加入 $limit 限制輸出結果
// 		pipeline = append(pipeline, bson.M{"$limit": opts.Limit})
// 	}

// 	// 执行聚合查询
// 	cursor, err := h.read.GetDB().Collection(db.ColName_UserRecord).Aggregate(ctx, pipeline)
// 	if err != nil {
// 		return nil, 0, errorcode.Server(err)
// 	}
// 	defer cursor.Close(ctx)

// 	// 解码结果
// 	var results []*db.UserRecord

// 	if err := cursor.All(ctx, &results); err != nil {
// 		return nil, 0, errorcode.Server(err)
// 	}

// 	return results, totalCount.Count, nil
// }

// // 創建指定日期 UserRecord_日期
// // datas 只能放相同日期的資料
// func (h *Handler) CreateUserRecord(ctx context.Context, data []*db.UserRecord) *errorcode.Error {
// 	col := h.write.GetDB().Collection(db.ColName_UserRecord + mtime.GetTimeFormatUnix(data[0].CreateTime, mtime.Format_YMD))

// 	// 只有一筆資料
// 	if len(data) == 1 {
// 		if _, err := col.InsertOne(ctx, data[0]); err != nil {
// 			return errorcode.Server(err)
// 		}
// 		return nil
// 	}

// 	// 多筆資料
// 	// 转换为 []interface{}
// 	iData := make([]interface{}, len(data))
// 	for i, v := range data {
// 		iData[i] = v
// 	}

// 	if _, err := col.InsertMany(ctx, iData); err != nil {
// 		return errorcode.Server(err)
// 	}

// 	return nil
// }

// // 取得 指定日期的 ResultBalance總和
// // map[2024_04_16:-7 2024_04_17:50]
// // func (h *Handler) GetUserRecordTotalResultBalance(ctx context.Context, filter *db.UserRecordOpt) (map[string]int64, *errorcode.Error) {
// func (h *Handler) GetUserRecordTotalResultBalance(ctx context.Context, filter *db.UserRecordOpt) ([]*db.UserRecordTotalResult, *errorcode.Error) {
// 	// 取得日期範圍
// 	dates := mtime.GetDateRange(filter.StartTimeUnix, filter.EndTimeUnix)

// 	// 构建聚合管道
// 	pipeline := make([]bson.M, len(dates)+2)

// 	for i := range dates {
// 		pipeline[i] = bson.M{
// 			"$unionWith": bson.M{
// 				"coll": db.ColName_UserRecord + dates[i],
// 				"pipeline": []bson.M{
// 					{"$set": bson.M{"date": dates[i]}},
// 				},
// 			},
// 		}
// 	}

// 	pipeline[len(pipeline)-2] = bson.M{"$match": filter.Filter_Mgo()}

// 	pipeline[len(pipeline)-1] = bson.M{
// 		"$group": bson.M{
// 			"_id":   "$date", // unionWith > pipeline > set date
// 			"total": bson.M{"$sum": "$result_balance"},
// 		},
// 	}

// 	// 执行聚合查询
// 	cursor, err := h.read.GetDB().Collection(db.ColName_UserRecord).Aggregate(ctx, pipeline)
// 	if err != nil {
// 		return nil, errorcode.Server(err)
// 	}
// 	defer cursor.Close(ctx)

// 	// 解码结果
// 	var results []bson.M
// 	if err := cursor.All(ctx, &results); err != nil {
// 		return nil, errorcode.Server(err)
// 	}

// 	// 格式化输出
// 	// output := make(map[string]int64)
// 	// for _, result := range results {
// 	// 	output[result["_id"].(string)] = result["total"].(int64)
// 	// }

// 	// 格式化输出
// 	output := make([]*db.UserRecordTotalResult, len(results))
// 	for i, result := range results {
// 		output[i] = &db.UserRecordTotalResult{
// 			Date:  result["_id"].(string),
// 			Total: result["total"].(int64),
// 		}
// 	}

// 	return output, nil
// }
