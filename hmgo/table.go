package hmgo

// import (
// 	"context"

// 	"github.com/chu16537/module_master/errorcode"
// 	"github.com/chu16537/module_master/proto/db"
// )

// // 取得 table
// func (h *Handler) GetTable(ctx context.Context, filter *db.TableOpt, opts *db.FindOpt) ([]*db.Table, *errorcode.Error) {
// 	cur, err := h.find(ctx, db.ColName_Table, filter.Filter_Mgo(), opts)
// 	if err != nil {
// 		return nil, errorcode.Server(err)
// 	}

// 	defer cur.Close(ctx)

// 	var results []*db.Table

// 	if err := cur.All(ctx, &results); err != nil {
// 		return nil, errorcode.Server(err)
// 	}

// 	return results, nil
// }

// // 取得count Table
// func (h *Handler) GetTableCount(ctx context.Context, filter *db.TableOpt) (int64, *errorcode.Error) {
// 	count, err := h.count(ctx, db.ColName_Table, filter.Filter_Mgo())
// 	if err != nil {
// 		return 0, errorcode.Server(err)
// 	}
// 	return count, nil
// }

// // 創建 table
// func (h *Handler) CreateTable(ctx context.Context, data *db.Table) *errorcode.Error {
// 	// 取得id
// 	tableID, errC := h.incr(ctx, db.ColName_Table)
// 	if errC != nil {
// 		return errC
// 	}

// 	data.ID = uint64(tableID)

// 	err := h.create(ctx, db.ColName_Table, data)
// 	if err != nil {
// 		return errorcode.Server(err)
// 	}

// 	return nil
// }

// // 更新 table
// func (h *Handler) UpdateTable(ctx context.Context, filter *db.TableOpt, data map[string]interface{}) (int64, *errorcode.Error) {
// 	c, err := h.update(ctx, db.ColName_Table, filter.Filter_Mgo(), data)
// 	if err != nil {
// 		return c, errorcode.Server(err)
// 	}

// 	return c, nil
// }
