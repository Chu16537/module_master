package mmgo

import (
	"context"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/proto/db"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IFiltetOpt interface {
	ToMgo() bson.M
}

// 創建資料表
// uniqueKeys 可以多個索引為一組 (ex  user_id, table_id) 每一組是一組[]string
// ex []string{[]string{"user_id", "table_id"},[]string{"user_id"}}
func (h *Handler) CreateCollection(ctx context.Context, colName string, uniqueKeys [][]string) *errorcode.Error {
	col := h.write.db.Collection(colName)

	for _, keys := range uniqueKeys {

		// 記錄keys
		ks := bson.D{}
		for _, k := range keys {
			ks = append(ks, bson.E{Key: k, Value: 1})
		}

		// 創建索引
		model := mongo.IndexModel{
			Keys:    ks,
			Options: options.Index().SetUnique(true),
		}

		_, err := col.Indexes().CreateOne(ctx, model)
		if err != nil {
			return errorcode.New(errorcode.Code_DB_Insert_Error, err)
		}
	}

	return nil
}

// 請傳入 *T
func FindOne[T any](h *Handler, ctx context.Context, colName string, filter IFiltetOpt) (T, *errorcode.Error) {
	col := h.read.db.Collection(colName)

	var result T
	err := col.FindOne(ctx, filter.ToMgo()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, nil // 如果沒找到，回傳 nil 而不是錯誤
		}
		return result, errorcode.New(errorcode.Code_DB_Find_Error, errors.Wrap(err, err.Error()))
	}
	return result, nil
}

func Find[T any](h *Handler, ctx context.Context, colName string, filter IFiltetOpt, findOpt *db.FindOpt) ([]T, *errorcode.Error) {
	col := h.read.db.Collection(colName)
	o := findOpt.ToMgo()

	cursor, err := col.Find(ctx, filter.ToMgo(), o)
	if err != nil {
		return nil, errorcode.New(errorcode.Code_DB_Find_Error, errors.Wrap(err, err.Error()))
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, errorcode.New(errorcode.Code_DB_Cursor_Error, errors.Wrap(err, err.Error()))
	}

	return results, nil
}

func Count(h *Handler, ctx context.Context, name string, filter IFiltetOpt) (int64, *errorcode.Error) {
	count, err := h.read.db.Collection(name).CountDocuments(ctx, filter.ToMgo())
	if err != nil {
		return 0, errorcode.New(errorcode.Code_DB_Count_Error, errors.Wrap(err, err.Error()))
	}
	return count, nil
}

func Insert(h *Handler, ctx context.Context, name string, data interface{}) *errorcode.Error {
	col := h.write.db.Collection(name)
	_, err := col.InsertOne(ctx, data)
	if err != nil {
		return errorcode.New(errorcode.Code_DB_Insert_Error, errors.Wrap(err, err.Error()))
	}
	return nil
}

func InsertMany(h *Handler, ctx context.Context, name string, data ...interface{}) *errorcode.Error {
	col := h.write.db.Collection(name)
	_, err := col.InsertMany(ctx, data)
	if err != nil {
		return errorcode.New(errorcode.Code_DB_InsertMany_Error, errors.Wrap(err, err.Error()))
	}

	return nil
}

func Update(h *Handler, ctx context.Context, name string, filter IFiltetOpt, data map[string]interface{}) (int64, *errorcode.Error) {
	col := h.write.db.Collection(name)

	update := bson.M{"$set": data}

	r, err := col.UpdateMany(ctx, filter.ToMgo(), update)
	if err != nil {
		return 0, errorcode.New(errorcode.Code_DB_Update_Error, errors.Wrap(err, err.Error()))
	}

	// MatchedCount 符合條件的資料數量
	// ModifiedCount 更改的資料數量
	return r.ModifiedCount, nil
}

func Delete(h *Handler, ctx context.Context, name string, filter bson.M) (int64, *errorcode.Error) {
	col := h.write.db.Collection(name)

	r, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return 0, errorcode.New(errorcode.Code_DB_Delete_Error, errors.Wrap(err, err.Error()))
	}

	// DeletedCount 删除的資料數量
	return r.DeletedCount, nil
}

func Aggregate[T any](h *Handler, ctx context.Context, colName string, pipeline interface{}, opts *db.FindOpt) ([]T, *errorcode.Error) {
	col := h.read.db.Collection(colName)

	cursor, err := col.Aggregate(ctx, pipeline, opts.ToAggregate())
	if err != nil {
		return nil, errorcode.New(errorcode.Code_DB_Aggregate_Error, errors.Wrap(err, err.Error()))
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, errorcode.New(errorcode.Code_DB_Cursor_Error, errors.Wrap(err, err.Error()))
	}

	return results, nil
}

// 批量执行多个写操作（包括插入、更新和删除）。
func BulkWrite(h *Handler, ctx context.Context, name string, operations []mongo.WriteModel) (*mongo.BulkWriteResult, *errorcode.Error) {
	col := h.write.db.Collection(name)

	result, err := col.BulkWrite(ctx, operations)
	if err != nil {
		return nil, errorcode.New(errorcode.Code_DB_BulkWrite_Error, errors.Wrap(err, err.Error()))
	}

	return result, nil
}

// 執行 Transaction
func Transaction(h *Handler, ctx context.Context, callback func(sctx mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) *errorcode.Error {
	// 強制使用 讀讀 這樣 查詢跟寫入都是使用讀db
	nh := h.GetWW()

	s, err := nh.write.client.StartSession()
	if err != nil {
		return errorcode.New(errorcode.Code_DB_Transaction_Error, errors.Wrap(err, err.Error()))
	}

	// 無論事務是否成功，我們都必須結束會話以釋放資源。
	defer s.EndSession(ctx)

	// 執行 transaction
	if _, err := s.WithTransaction(ctx, callback); err != nil {
		return errorcode.New(errorcode.Code_DB_Transaction_Error, errors.Wrap(err, err.Error()))
	}

	return nil
}

// 自增id
func Incr(h *Handler, ctx context.Context, tagColName string) (int, *errorcode.Error) {
	col := h.write.db.Collection(db.ColName_Counters)

	filter := bson.M{"col_name": tagColName}
	update := bson.M{"$inc": bson.M{"value": 1}}
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	counter := &db.Counter{}
	err := col.FindOneAndUpdate(ctx, filter, update, options).Decode(counter)

	if err != nil {
		return 0, errorcode.New(errorcode.Code_DB_Update_Error, errors.New(err.Error()))
	}

	return counter.Value, nil
}

// 自增id多筆
func Incrs(h *Handler, ctx context.Context, tagColName string, count int) ([]int, *errorcode.Error) {
	col := h.write.db.Collection(db.ColName_Counters)

	filter := bson.M{"_id": tagColName}
	update := bson.M{"$inc": bson.M{"value": count}}
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	counter := &db.Counter{}
	err := col.FindOneAndUpdate(ctx, filter, update, options).Decode(counter)

	if err != nil {
		return []int{}, errorcode.New(errorcode.Code_DB_Update_Error, errors.Wrap(err, err.Error()))
	}

	result := []int{counter.Value - count + 1, counter.Value}

	return result, nil
}
