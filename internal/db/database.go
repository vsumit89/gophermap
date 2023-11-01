package db

import "gophermap/internal/db/postgres"

type IDatabase interface {
	Connect() error
	Disconnect() error
	GetClient() interface{}
}

func NewDBService() IDatabase {
	return postgres.NewPostgresAdapter()
}
