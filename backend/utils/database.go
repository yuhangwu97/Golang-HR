package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func WithDBContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return context.WithTimeout(context.Background(), timeout)
}

func ValidateObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

func SetTimestamps(entity interface{}, isUpdate bool) {
	now := time.Now()
	
	switch v := entity.(type) {
	case map[string]interface{}:
		if !isUpdate {
			v["created_at"] = now
		}
		v["updated_at"] = now
	}
}

func CreateUpdateSet(fields map[string]interface{}) bson.M {
	SetTimestamps(fields, true)
	return bson.M{"$set": fields}
}

func CreateFilter(field, value string) bson.M {
	if field == "_id" || field == "id" {
		if objectID, err := primitive.ObjectIDFromHex(value); err == nil {
			return bson.M{"_id": objectID}
		}
	}
	return bson.M{field: value}
}