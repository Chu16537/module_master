package hmgo

import (
	"context"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// 事務
// 預設失敗重新執行5次
func (h *Handler) startTransaction(ctx context.Context, fn func(h *Handler, sessionContext mongo.SessionContext) (interface{}, *errorcode.Error)) (err *errorcode.Error) {
	trans := h.trans()

	session, e := trans.write.GetClient().StartSession()
	if e != nil {
		return errorcode.Server(e)
	}

	fu := func(sessionContext mongo.SessionContext) (interface{}, error) {
		_, e := fn(trans, sessionContext)

		if e != nil {
			err = e
			return nil, e.Err()
		}

		err = nil
		return nil, nil
	}

	defer func() {
		// 在捕获异常时执行回滚
		if r := recover(); r != nil {
			// 如果有 panic，设置 err 以便在函数末尾返回错误
			err = errorcode.Server(errors.Errorf("panic occurred: %v", r))
			return
		}
		// 处理事务的正常结束或异常结束
		if err != nil {
			session.AbortTransaction(ctx)
		} else {
			// 提交事务
			e := session.CommitTransaction(ctx)
			if e != nil {
				session.AbortTransaction(ctx)
				err = errorcode.Server(e)
			}
		}
		// 無論事務是否成功，我們都必須結束會話以釋放資源。
		session.EndSession(ctx)
	}()

	// 执行事务操作
	_, e = session.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
		return fu(sessionContext)
	})

	if err != nil {
		return
	}

	if e != nil {
		// 執行完成fn
		err = errorcode.Server(e)
		return
	}

	// 執行完成 WithTransaction
	return nil
}
