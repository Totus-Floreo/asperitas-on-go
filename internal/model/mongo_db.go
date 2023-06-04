package model

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongoDB interface {
	Collection(string, ...*options.CollectionOptions) ICollection
}

type MyMongoDatabase struct {
	*mongo.Database
}

func (db MyMongoDatabase) Collection(name string, opts ...*options.CollectionOptions) ICollection {
	coll := MyMongoCollection{
		db.Database.Collection(name, opts...),
	}
	return coll
}
