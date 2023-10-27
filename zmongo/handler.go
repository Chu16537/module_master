package zmongo

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *Handler) Find(colName string, filter bson.M, opts *options.FindOptions, obj interface{}) ([]interface{}, error) {
	col := h.db.Collection(colName)
	cur, err := col.Find(h.ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	defer cur.Close(h.ctx)

	objectType := reflect.TypeOf(obj).Elem()
	var list = make([]interface{}, 0)

	for cur.Next(context.Background()) {
		result := reflect.New(objectType).Interface()
		err := cur.Decode(result)

		if err != nil {
			return nil, err
		}

		list = append(list, result)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
