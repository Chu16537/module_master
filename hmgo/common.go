package hmgo

// import (
// 	"context"

// 	"github.com/chu16537/module_master/proto/db"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func (h *Handler) find(ctx context.Context, name string, filter bson.M, opts *db.FindOpt) (*mongo.Cursor, error) {
// 	col := h.read.GetDB().Collection(name)

// 	o := options.Find()

// 	if opts != nil {
// 		opts.ToMgo()
// 		// 时间排序
// 		// o := o.SetSort(bson.M{"create_time": 1})
// 		o = o.SetSkip(int64(opts.Start)).SetLimit(int64(opts.Limit))
// 	}

// 	return col.Find(ctx, filter, o)
// }

// func (h *Handler) count(ctx context.Context, name string, filter bson.M) (int64, error) {
// 	return h.read.GetDB().Collection(name).CountDocuments(ctx, filter)
// }

// func (h *Handler) create(ctx context.Context, name string, data interface{}) error {
// 	col := h.write.GetDB().Collection(name)
// 	_, err := col.InsertOne(ctx, data)
// 	return err
// }

// func (h *Handler) update(ctx context.Context, name string, filter bson.M, data map[string]interface{}) (int64, error) {
// 	col := h.write.GetDB().Collection(name)

// 	update := bson.M{"$set": data}

// 	r, err := col.UpdateMany(ctx, filter, update)
// 	if err != nil {
// 		return 0, err
// 	}

// 	// MatchedCount 符合條件的資料數量
// 	// ModifiedCount 更改的資料數量
// 	// fmt.Println(r.MatchedCount, r.ModifiedCount, r.UpsertedCount, r.UpsertedID)
// 	return r.ModifiedCount, nil
// }
