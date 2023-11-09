package zmongo

// IndexesCreateOne
type CreateIndex struct {
	ColName string
	Key     string
}

// 計數器
type Counter struct {
	ID    string `bson:"_id"`
	Value int    `bson:"value"`
}
