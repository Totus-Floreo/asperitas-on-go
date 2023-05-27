package mongo_repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBReadersPool struct {
	connections chan *mongo.Client
}

func NewDBReadersPool(uri string, maxConn int) (*DBReadersPool, error) {
	pool := &DBReadersPool{
		connections: make(chan *mongo.Client, maxConn),
	}

	for i := 0; i < maxConn; i++ {
		client, err := mongo.NewClient(options.Client().ApplyURI(uri))
		if err != nil {
			return nil, err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = client.Connect(ctx)
		if err != nil {
			return nil, err
		}

		pool.connections <- client
	}

	return pool, nil
}

func (p *DBReadersPool) GetConnection() *mongo.Client {
	return <-p.connections
}

func (p *DBReadersPool) ReleaseConnection(conn *mongo.Client) {
	p.connections <- conn
}
