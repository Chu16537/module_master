package hmgo

import (
	"context"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/proto/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 自增id
func (h *Handler) incr(ctx context.Context, tagColName string) (int, *errorcode.Error) {
	col := h.write.GetDB().Collection(db.ColName_Counters)

	filter := bson.M{"_id": tagColName}
	update := bson.M{"$inc": bson.M{"value": 1}}
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	counter := &db.Counter{}
	err := col.FindOneAndUpdate(ctx, filter, update, options).Decode(counter)

	if err != nil {
		return 0, errorcode.Server(err)
	}

	return counter.Value, nil
}

// 自增id多筆
func (h *Handler) incrs(ctx context.Context, tagColName string, count int) ([]int, *errorcode.Error) {
	col := h.write.GetDB().Collection(db.ColName_Counters)

	filter := bson.M{"_id": tagColName}
	update := bson.M{"$inc": bson.M{"value": count}}
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	counter := &db.Counter{}
	err := col.FindOneAndUpdate(ctx, filter, update, options).Decode(counter)

	if err != nil {
		return nil, errorcode.Server(err)
	}

	result := []int{counter.Value - count + 1, counter.Value}

	return result, nil
}
