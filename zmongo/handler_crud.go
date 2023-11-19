package zmongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 查一筆資料
func (h *Handler) FindOne(colName string, obj interface{}, filter bson.M, opts ...*options.FindOneOptions) error {
	col := h.db.Collection(colName)

	err := col.FindOne(h.ctx, filter, opts...).Decode(obj)
	if err != nil {
		return err
	}

	return nil
}

// 查多筆資料
func (h *Handler) FindArray(colName string, filter bson.M, opts ...*options.FindOptions) ([]interface{}, error) {
	col := h.db.Collection(colName)

	cur, err := col.Find(h.ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	defer cur.Close(h.ctx)

	var list = make([]interface{}, 0)

	// 使用 cursor.All 一次性获取所有结果
	err = cur.All(h.ctx, &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// 新增一筆資料
func (h *Handler) InsertOne(colName string, obj interface{}, opts ...*options.InsertOneOptions) error {
	col := h.db.Collection(colName)

	// 插入数据
	_, err := col.InsertOne(h.ctx, obj, opts...)

	// 錯誤不是重複
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}

	return nil
}

// 刪除一筆資料
func (h *Handler) DelOne(colName string, filter bson.M, opts ...*options.DeleteOptions) error {
	col := h.db.Collection(colName)

	_, err := col.DeleteOne(h.ctx, filter, opts...)

	if err != nil {
		return err
	}

	// 反正結果是要刪除 沒有刪除 代表裡面應該也沒這裡資料
	// if result.DeletedCount == 0 {
	// 	return fmt.Errorf("DelOne not data")
	// }

	return nil
}

// 更新一筆資料
func (h *Handler) UpdateOne(colName string, filter bson.M, update bson.M, opts ...*options.UpdateOptions) error {
	col := h.db.Collection(colName)

	// 更新数据
	result, err := col.UpdateOne(h.ctx, filter, update, opts...)

	if err != nil {
		return err
	}

	// 沒有資料更新到
	if result.MatchedCount == 0 {
		return fmt.Errorf("UpdateOne: no document matched")
	}

	return nil
}

// 更新多筆資料 每個值都相同
func (h *Handler) UpdateManySameValue(colName string, filter bson.M, update bson.M, opts ...*options.UpdateOptions) error {
	col := h.db.Collection(colName)

	// 更新数据
	_, err := col.UpdateMany(h.ctx, filter, update, opts...)

	if err != nil {
		return err
	}

	return nil
}

// 更新多筆資料 每個數值都不一樣
func (h *Handler) UpdateManyDifferentValue(colName string, wm []mongo.WriteModel, opts ...*options.BulkWriteOptions) error {
	col := h.db.Collection(colName)

	_, err := col.BulkWrite(h.ctx, wm, opts...)

	if err != nil {
		return err
	}

	return nil
}
