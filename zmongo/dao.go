package zmongo

type Counter struct {
	ID    string `bson:"_id"`
	Value int    `bson:"value"`
}
