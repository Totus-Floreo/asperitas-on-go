package model

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ICollection interface {
	Aggregate(context.Context, interface{}, ...*options.AggregateOptions) (ICursor, error)
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult
	FindOneAndDelete(context.Context, interface{}, ...*options.FindOneAndDeleteOptions) *mongo.SingleResult
	DeleteOne(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateByID(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type MyMongoCollection struct {
	*mongo.Collection
}

func (coll MyMongoCollection) Aggregate(ctx context.Context, filter interface{}, opts ...*options.AggregateOptions) (ICursor, error) {
	return coll.Collection.Aggregate(ctx, filter, opts...)
}

func (coll MyMongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return coll.Collection.InsertOne(ctx, document, opts...)
}

func (coll MyMongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return coll.Collection.FindOne(ctx, filter, opts...)
}

func (coll MyMongoCollection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	return coll.Collection.FindOneAndDelete(ctx, filter, opts...)
}

func (coll MyMongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return coll.Collection.DeleteOne(ctx, filter, opts...)
}

func (coll MyMongoCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return coll.Collection.DeleteMany(ctx, filter, opts...)
}

func (coll MyMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return coll.Collection.UpdateOne(ctx, filter, update, opts...)
}

func (coll MyMongoCollection) UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return coll.Collection.UpdateByID(ctx, id, update, opts...)
}
