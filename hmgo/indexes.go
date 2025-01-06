package hmgo

// import (
// 	"context"

// 	"github.com/Chu16537/module_master/errorcode"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// // 索引
// type CreateIndex struct {
// 	ColName string
// 	Val     mongo.IndexModel
// }

// // 創建唯一索引
// func (h *Handler) CreateIndexs(ctx context.Context, cis []*CreateIndex) *errorcode.Error {
// 	for _, v := range cis {
// 		col := h.write.GetDB().Collection(v.ColName)

// 		// 创建索引
// 		if _, err := col.Indexes().CreateOne(ctx, v.Val); err != nil {
// 			return errorcode.Server(err)
// 		}
// 	}
// 	return nil
// }
