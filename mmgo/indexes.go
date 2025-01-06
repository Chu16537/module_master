package mmgo

import (
	"context"

	"github.com/Chu16537/module_master/errorcode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	return errorcode.Success()
}
