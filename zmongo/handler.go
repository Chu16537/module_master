package zmongo

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 創建唯一索引
func (h *Handler) IndexesCreateOne(ctx context.Context, colName string, key string) error {
	col := h.db.Collection(colName)

	im := mongo.IndexModel{
		Keys: bson.M{
			key: 1,
		},
		Options: options.Index().SetUnique(true),
	}
	// 创建索引
	if _, err := col.Indexes().CreateOne(ctx, im); err != nil {
		return err
	}

	return nil
}

// 查一筆資料
func (h *Handler) FindOne(ctx context.Context, colName string, filter bson.M, opts *options.FindOneOptions, obj interface{}) (interface{}, error) {
	col := h.db.Collection(colName)
	objectType := reflect.TypeOf(obj).Elem()
	err := col.FindOne(ctx, filter, opts).Decode(obj)
	if err != nil {
		return nil, err
	}
	return objectType, nil
}

// 查多筆資料
func (h *Handler) FindArray(ctx context.Context, colName string, filter bson.M, opts *options.FindOptions, obj interface{}) ([]interface{}, error) {
	col := h.db.Collection(colName)
	cur, err := col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	objectType := reflect.TypeOf(obj).Elem()
	var list = make([]interface{}, 0)

	for cur.Next(context.Background()) {
		result := reflect.New(objectType).Interface()
		err := cur.Decode(result)

		if err != nil {
			return nil, err
		}

		list = append(list, result)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

// 新增一筆資料 併發不重複
func (h *Handler) InsertOne(ctx context.Context, colName string, obj interface{}) error {
	col := h.db.Collection(colName)

	// 插入数据
	_, err := col.InsertOne(ctx, obj)
	// 錯誤不是重複
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}
	return nil
}

// 自增id
func (h *Handler) IncrementID(ctx context.Context, colName string, tagColName string) (int, error) {
	col := h.db.Collection(colName)

	filter := bson.D{{"_id", tagColName}}
	update := bson.M{"$inc": bson.M{"value": 1}}
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	res := col.FindOneAndUpdate(ctx, filter, update, options)

	if res.Err() != nil {
		return 0, res.Err()
	}

	var counter Counter
	if err := res.Decode(&counter); err != nil {
		return 0, err
	}

	return counter.Value, nil
}

// 事務
func (h *Handler) StartTransaction(ctx context.Context, cancel context.CancelFunc, f func(sessionContext mongo.SessionContext) (interface{}, error)) (interface{}, error) {

	defer cancel()

	session, err := h.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	// 开启事务
	result, err := session.WithTransaction(ctx, f)
	if err != nil {
		// 回滚事务
		errAbort := session.AbortTransaction(ctx)
		if errAbort != nil {
			fmt.Println("Error rolling back transaction:", errAbort)
		}
		return nil, err
	}

	return result, nil
}

// 更新多筆資料 每個數值都不一樣
func (h *Handler) BulkWrite(ctx context.Context, colName string, wm []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	col := h.db.Collection(colName)
	return col.BulkWrite(ctx, wm)
}
