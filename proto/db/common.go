package db

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FindOpt struct {
	Start uint64
	Limit uint64
	Mgo   *MongoFindOpt
}

type MongoFindOpt struct {
	Sort   map[string]int // 排序 // bson.D{{"price", 1}, {"name", -1}}：先按 price 升序，再按 name 降序。
	Fields map[string]int // 顯示得欄位 // bson.M{"name": 1, "age": 1}：只返回 name 和 age bson.M{"password": 0}：排除 password 欄位（0 表示不返回該欄位）。
}

func (o *FindOpt) ToMgo() *options.FindOptions {
	fo := options.Find()

	if o == nil {
		return fo
	}

	if o.Start > 0 {
		fo.SetSkip(int64(o.Start))
	}

	if o.Limit <= 0 {
		o.Limit = 10
	}

	if o.Limit > 0 {
		fo.SetLimit(int64(o.Limit))
	}

	if o.Mgo != nil {
		if len(o.Mgo.Sort) > 0 {
			fo.SetSort(o.Mgo.Sort)
		}
		if len(o.Mgo.Fields) > 0 {
			fo.SetProjection(o.Mgo.Fields)
		}
	}

	return fo
}

func (o *FindOpt) ToAggregate() *options.AggregateOptions {

	return nil
}

type TotalCount struct {
	Count int64 `bson:"count"`
}
