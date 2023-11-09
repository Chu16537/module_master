package zmongo

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 創建唯一索引
func (h *Handler) IndexesCreateOne(colName string, key string) error {
	col := h.db.Collection(colName)

	im := mongo.IndexModel{
		Keys: bson.M{
			key: 1,
		},
		Options: options.Index().SetUnique(true),
	}
	// 创建索引
	if _, err := col.Indexes().CreateOne(h.ctx, im); err != nil {
		return err
	}

	return nil
}

// 查一筆資料
func (h *Handler) FindOne(colName string, filter bson.M, opts *options.FindOneOptions, obj interface{}) (interface{}, error) {
	col := h.db.Collection(colName)
	objectType := reflect.TypeOf(obj).Elem()
	err := col.FindOne(h.ctx, filter, opts).Decode(obj)
	if err != nil {
		return nil, err
	}
	return objectType, nil
}

// 查多筆資料
func (h *Handler) FindArray(colName string, filter bson.M, opts *options.FindOptions, obj interface{}) ([]interface{}, error) {
	col := h.db.Collection(colName)
	cur, err := col.Find(h.ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	defer cur.Close(h.ctx)

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
func (h *Handler) InsertOne(colName string, obj interface{}) error {
	col := h.db.Collection(colName)

	// 插入数据
	_, err := col.InsertOne(h.ctx, obj)
	// 錯誤不是重複
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}
	return nil
}

// 自增id
func (h *Handler) IncrementID(colName string, tagColName string) (int, error) {
	col := h.db.Collection(colName)

	filter := bson.D{{"_id", tagColName}}
	update := bson.M{"$inc": bson.M{"value": 1}}
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	res := col.FindOneAndUpdate(h.ctx, filter, update, options)

	if res.Err() != nil {
		return 0, res.Err()
	}

	var counter Counter
	if err := res.Decode(&counter); err != nil {
		return 0, err
	}

	return counter.Value, nil
}
