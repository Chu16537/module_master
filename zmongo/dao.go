package zmongo

import "go.mongodb.org/mongo-driver/mongo"

// IndexesCreateOne
type CreateIndex struct {
	ColName string
	Val     mongo.IndexModel
}

// 計數器
type Counter struct {
	ID    string `bson:"_id"`
	Value int    `bson:"value"`
}
