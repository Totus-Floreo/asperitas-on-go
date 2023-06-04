package mongo_repository

import (
	"context"
	"time"

	"github.com/Totus-Floreo/asperitas-on-go/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBReadersPool struct {
	connections chan model.IClient
}

func NewDBReadersPool(uri string, maxConn int) (*DBReadersPool, error) {
	pool := &DBReadersPool{
		connections: make(chan model.IClient, maxConn),
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

		myClient := model.MyMongoClient{
			Client: client,
		}

		pool.connections <- myClient
	}

	return pool, nil
}

func (p *DBReadersPool) GetConnection() model.IClient {
	return <-p.connections
}

func (p *DBReadersPool) ReleaseConnection(conn model.IClient) {
	p.connections <- conn
}
