package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

func NewBaseRepository(collection *mongo.Collection) *BaseRepository {
	return &BaseRepository{
		collection: collection,
		timeout:    5 * time.Second,
	}
}

func (r *BaseRepository) WithTimeout(timeout time.Duration) *BaseRepository {
	r.timeout = timeout
	return r
}

func (r *BaseRepository) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), r.timeout)
}

func (r *BaseRepository) Create(entity interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := r.getContext()
	defer cancel()
	
	return r.collection.InsertOne(ctx, entity)
}

func (r *BaseRepository) FindByID(id string, result interface{}) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	ctx, cancel := r.getContext()
	defer cancel()
	
	return r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(result)
}

func (r *BaseRepository) Find(filter bson.M, result interface{}) error {
	ctx, cancel := r.getContext()
	defer cancel()
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	
	return cursor.All(ctx, result)
}

func (r *BaseRepository) UpdateByID(id string, update bson.M) (*mongo.UpdateResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	ctx, cancel := r.getContext()
	defer cancel()
	
	return r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
}

func (r *BaseRepository) DeleteByID(id string) (*mongo.DeleteResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	ctx, cancel := r.getContext()
	defer cancel()
	
	return r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
}

func (r *BaseRepository) Count(filter bson.M) (int64, error) {
	ctx, cancel := r.getContext()
	defer cancel()
	
	return r.collection.CountDocuments(ctx, filter)
}

func (r *BaseRepository) FindWithPagination(filter bson.M, page, pageSize int, result interface{}) error {
	ctx, cancel := r.getContext()
	defer cancel()
	
	skip := (page - 1) * pageSize
	
	opts := options.Find()
	opts.SetSkip(int64(skip))
	opts.SetLimit(int64(pageSize))
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	
	return cursor.All(ctx, result)
}

func (r *BaseRepository) Exists(filter bson.M) (bool, error) {
	ctx, cancel := r.getContext()
	defer cancel()
	
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

func (r *BaseRepository) FindOne(filter bson.M, result interface{}) error {
	ctx, cancel := r.getContext()
	defer cancel()
	
	return r.collection.FindOne(ctx, filter).Decode(result)
}