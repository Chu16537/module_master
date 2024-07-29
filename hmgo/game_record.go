package hmgo

import (
	"context"

	"github.com/Chu16537/gomodule/errorcode"
	"github.com/Chu16537/gomodule/mtime"
	"github.com/Chu16537/gomodule/proto/db"
	"go.mongodb.org/mongo-driver/bson"
)

// 取得 GameRecord
func (h *Handler) GetGameRecord(ctx context.Context, filter *db.GameRecordOpt, opts *db.FindOpt) ([]*db.GameRecord, int64, *errorcode.Error) {
	dates := mtime.GetDateRange(filter.StartTimeUnix, filter.EndTimeUnix)

	pipeline := make([]bson.M, len(dates)+2)

	for i := range dates {
		pipeline[i] = bson.M{
			"$unionWith": bson.M{
				"coll": db.ColName_GameRecord + dates[i],
			},
		}
	}

	pipeline[len(pipeline)-2] = bson.M{
		"$match": filter.Filter_Mgo(),
	}

	pipeline[len(pipeline)-1] = bson.M{"$count": "count"}

	// 执行聚合查询以获取总数
	totalCountCursor, err := h.read.GetDB().Collection(db.ColName_GameRecord).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, errorcode.Server(err)
	}
	defer totalCountCursor.Close(ctx)

	totalCount := &db.TotalCount{}
	if totalCountCursor.Next(ctx) {
		if err := totalCountCursor.Decode(totalCount); err != nil {
			return nil, 0, errorcode.Server(err)
		}
	}

	pipeline = pipeline[:len(pipeline)-1]

	if opts != nil {
		opts.ToMgo()
		// 加入 $skip 跳過前面的紀錄
		pipeline = append(pipeline, bson.M{"$skip": opts.Start})
		// 加入 $limit 限制輸出結果
		pipeline = append(pipeline, bson.M{"$limit": opts.Limit})
	}

	// 执行聚合查询
	cursor, err := h.read.GetDB().Collection(db.ColName_GameRecord).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, errorcode.Server(err)
	}
	defer cursor.Close(ctx)

	// 解码结果
	var results []*db.GameRecord
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, errorcode.Server(err)
	}

	return results, totalCount.Count, nil
}

// 創建 GameRecord
func (h *Handler) createGameRecord(ctx context.Context, gr *db.GameRecord) *errorcode.Error {
	col := h.write.GetDB().Collection(db.ColName_GameRecord + mtime.GetTimeFormatUnix(gr.CreateTime, mtime.Format_YMD))

	if _, err := col.InsertOne(ctx, gr); err != nil {
		return errorcode.Server(err)
	}

	return nil
}
