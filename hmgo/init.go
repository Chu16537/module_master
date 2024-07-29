package hmgo

import (
	"context"

	"github.com/Chu16537/gomodule/errorcode"
	"github.com/Chu16537/gomodule/proto/db"
	"github.com/Chu16537/gomodule/ztime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
啟動執行一次

只要執行一次的即可
設定 資料表 IndexModel

測試:完成
*/
func (h *Handler) Init() *errorcode.Error {
	cis := []*CreateIndex{
		{
			ColName: db.ColName_Club,
			Val: mongo.IndexModel{
				Keys: bson.D{
					{Key: "id", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		{
			ColName: db.ColName_Club,
			Val: mongo.IndexModel{
				Keys: bson.D{
					{Key: "invitation_code", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		{
			ColName: db.ColName_Club_User_Info,
			Val: mongo.IndexModel{
				Keys: bson.D{
					{Key: "account", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		{
			ColName: db.ColName_Club_User_Info,
			Val: mongo.IndexModel{
				Keys: bson.D{
					{Key: "user_id", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		{
			ColName: db.ColName_Club_User_Info,
			Val: mongo.IndexModel{
				Keys: bson.D{
					{Key: "token", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		{
			ColName: db.ColName_Table,
			Val: mongo.IndexModel{
				Keys: bson.D{
					{Key: "id", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		{
			ColName: db.ColName_GameWallet,
			Val: mongo.IndexModel{
				Keys: bson.D{
					{Key: "user_id", Value: 1},
					{Key: "table_id", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
	}

	err := h.CreateIndexs(h.ctx, cis)
	if err != nil {
		return err
	}

	return err
}

/*
每日要創建的資料表

測試:完成
*/
func (h *Handler) EveryDay(ctx context.Context, unixTime int64, days int) *errorcode.Error {
	// 預設都先創建一個禮拜

	cis := []*CreateIndex{}
	for i := 0; i < days; i++ {
		date := ztime.GetTimeFormatUnix(unixTime+int64(ztime.Day_Sceond*i), ztime.Format_YMD)

		// wallet log
		wl := &CreateIndex{
			ColName: db.ColName_WalletLog + date,
			Val: mongo.IndexModel{
				Keys: bson.D{
					{Key: "order_id", Value: 1},
					{Key: "user_id", Value: 1},
					{Key: "event_type", Value: 1},
					{Key: "wallet_type", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		}

		cis = append(cis, wl)

		// game record
		gr := &CreateIndex{
			ColName: db.ColName_GameRecord + date,
			Val: mongo.IndexModel{
				Keys: bson.D{
					{Key: "game_record_id", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		}

		cis = append(cis, gr)

		// user record
		ur := &CreateIndex{
			ColName: db.ColName_UserRecord + date,
			Val: mongo.IndexModel{
				Keys: bson.D{
					{Key: "user_id", Value: 1},
					{Key: "game_record_id", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		}

		cis = append(cis, ur)
	}

	err := h.CreateIndexs(ctx, cis)
	if err != nil {
		return err
	}

	return nil
}
