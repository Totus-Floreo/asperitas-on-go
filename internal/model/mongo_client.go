package model

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IClient interface {
	Database(string, ...*options.DatabaseOptions) IMongoDB
}

type MyMongoClient struct {
	*mongo.Client
}

func (m MyMongoClient) Database(name string, opts ...*options.DatabaseOptions) IMongoDB {
	mongoDatabase := &MyMongoDatabase{
		m.Client.Database(name, opts...),
	}
	return mongoDatabase
}
