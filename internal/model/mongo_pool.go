package model

type IDBReadersPool interface {
	GetConnection() IClient
	ReleaseConnection(conn IClient)
}
