package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FindOpt struct {
	Start  uint64
	Limit  uint64
	Sort   bson.D
	Fields bson.M
}

func checkTimeUxin(start, end int64) (int64, int64) {
	// 時間 start 比 end 小
	if start > end {
		return end, start
	}
	return start, end
}

func (o *FindOpt) ToMgo() *options.FindOptions {
	fo := options.Find()

	if o == nil {
		return fo
	}

	if o.Start > 0 {
		fo.SetSkip(int64(o.Start))
	}
	if o.Limit > 0 {
		fo.SetLimit(int64(o.Limit))
	}
	if len(o.Sort) > 0 {
		fo.SetSort(o.Sort)
	}
	if len(o.Fields) > 0 {
		fo.SetProjection(o.Fields)
	}

	return fo
}

func (o *FindOpt) ToAggregate() *options.AggregateOptions {

	return nil
}

type TotalCount struct {
	Count int64 `bson:"count"`
}
