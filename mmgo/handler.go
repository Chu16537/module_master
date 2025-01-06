package mmgo

import (
	"context"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/proto/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Find[T any](h *Handler, ctx context.Context, colName string, filter bson.M, opts *db.FindOpt) ([]T, *errorcode.Error) {
	col := h.read.db.Collection(colName)
	o := opts.ToMgo()

	cursor, err := col.Find(ctx, filter, o)
	if err != nil {
		return nil, errorcode.New(errorcode.Code_DB_Find_Error, err)
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, errorcode.New(errorcode.Code_DB_Cursor_Error, err)
	}

	return results, errorcode.Success()
}

func Count(h *Handler, ctx context.Context, name string, filter bson.M) (int64, *errorcode.Error) {
	count, err := h.read.db.Collection(name).CountDocuments(ctx, filter)
	if err != nil {
		return 0, errorcode.New(errorcode.Code_DB_Count_Error, err)
	}
	return count, errorcode.Success()
}

func Create(h *Handler, ctx context.Context, name string, data interface{}) *errorcode.Error {
	col := h.write.db.Collection(name)
	_, err := col.InsertOne(ctx, data)
	if err != nil {
		return errorcode.New(errorcode.Code_DB_Insert_Error, err)
	}
	return errorcode.Success()
}

func CreateMany(h *Handler, ctx context.Context, name string, data []interface{}) *errorcode.Error {
	col := h.write.db.Collection(name)
	_, err := col.InsertMany(ctx, data)
	if err != nil {
		return errorcode.New(errorcode.Code_DB_InsertMany_Error, err)
	}
	return errorcode.Success()
}

// MatchedCount 符合條件的資料數量
// ModifiedCount 更改的資料數量
func Update(h *Handler, ctx context.Context, name string, filter bson.M, data map[string]interface{}) (int64, *errorcode.Error) {
	col := h.write.db.Collection(name)

	update := bson.M{"$set": data}

	r, err := col.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, errorcode.New(errorcode.Code_DB_Update_Error, err)
	}

	return r.ModifiedCount, errorcode.Success()
}

func Delete(h *Handler, ctx context.Context, name string, filter bson.M) (int64, *errorcode.Error) {
	col := h.write.db.Collection(name)

	r, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return 0, errorcode.New(errorcode.Code_DB_Delete_Error, err)
	}

	// DeletedCount 删除的資料數量
	return r.DeletedCount, errorcode.Success()
}

func Aggregate[T any](h *Handler, ctx context.Context, colName string, pipeline interface{}, opts *db.FindOpt) ([]T, *errorcode.Error) {
	col := h.read.db.Collection(colName)

	cursor, err := col.Aggregate(ctx, pipeline, opts.ToAggregate())
	if err != nil {
		return nil, errorcode.New(errorcode.Code_DB_Aggregate_Error, err)
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, errorcode.New(errorcode.Code_DB_Cursor_Error, err)
	}

	return results, errorcode.Success()
}

// 批量执行多个写操作（包括插入、更新和删除）。
func BulkWrite(h *Handler, ctx context.Context, name string, operations []mongo.WriteModel) (*mongo.BulkWriteResult, *errorcode.Error) {
	col := h.write.db.Collection(name)

	result, err := col.BulkWrite(ctx, operations)
	if err != nil {
		return nil, errorcode.New(errorcode.Code_DB_BulkWrite_Error, err)
	}

	return result, errorcode.Success()
}
