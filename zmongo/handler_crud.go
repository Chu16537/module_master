package zmongo

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 查一筆資料
func FindOne(db *mongo.Database, ctx context.Context, colName string, obj interface{}, filter bson.M, opts ...*options.FindOneOptions) error {
	col := db.Collection(colName)
	err := col.FindOne(ctx, filter, opts...).Decode(obj)
	if err != nil {
		return err
	}
	return nil
}

// 查多筆資料
func FindArray(db *mongo.Database, ctx context.Context, colName string, obj interface{}, filter bson.M, opts *options.FindOptions) ([]interface{}, error) {
	col := db.Collection(colName)
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

// 新增一筆資料
func InsertOne(db *mongo.Database, ctx context.Context, colName string, obj interface{}) error {
	col := db.Collection(colName)

	// 插入数据
	_, err := col.InsertOne(ctx, obj)

	// 錯誤不是重複
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}
	return nil
}

// 刪除一筆資料
func DelOne(db *mongo.Database, ctx context.Context, colName string, filter bson.M, opts *options.DeleteOptions) error {
	col := db.Collection(colName)

	// 更新数据
	_, err := col.DeleteOne(ctx, filter, opts)
	// 錯誤不是重複
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
func UpdateOne(db *mongo.Database, ctx context.Context, colName string, filter bson.M, update bson.M, opts *options.UpdateOptions) error {
	col := db.Collection(colName)

	// 更新数据
	result, err := col.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		return err
	}

	// 沒有資料更新到
	if result.MatchedCount == 0 {
		return fmt.Errorf("UpdateOne not data")
	}

	return nil
}

// 更新多筆資料 每個值都相同
func UpdateManySameValue(db *mongo.Database, ctx context.Context, colName string, filter bson.M, update bson.M, opts *options.UpdateOptions) error {
	col := db.Collection(colName)

	// 更新数据
	_, err := col.UpdateMany(ctx, filter, update, opts)

	if err != nil {
		return err
	}

	return nil
}

// 更新多筆資料 每個數值都不一樣
func UpdateManyDifferentValue(db *mongo.Database, ctx context.Context, colName string, wm []mongo.WriteModel) error {
	col := db.Collection(colName)

	_, err := col.BulkWrite(ctx, wm)

	if err != nil {
		return err
	}

	return nil
}
