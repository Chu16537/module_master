package zmongo

import (
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

// 自增id
func (h *Handler) IncrementID(colName string, tagColName string) (int, error) {
	col := h.db.Collection(colName)

	filter := bson.D{{"_id", tagColName}}
	update := bson.M{"$inc": bson.M{"value": 1}}
	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var counter Counter
	err := col.FindOneAndUpdate(h.ctx, filter, update, options).Decode(&counter)

	if err != nil {
		return 0, err
	}

	return counter.Value, nil
}

// 事務
func (h *Handler) StartTransaction(f func(sessionContext mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	session, err := h.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(h.ctx)

	// 开启事务
	result, err := session.WithTransaction(h.ctx, f)
	if err != nil {
		// 回滚事务
		errAbort := session.AbortTransaction(h.ctx)
		if errAbort != nil {
			return nil, errAbort
		}
		return nil, err
	}

	return result, nil
}
