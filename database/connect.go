package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(URI string) (*Database, error) {
	db := Database{}
	err := db.ConnectWithURI(URI)
	if err != nil {
		return nil, err
	}
	return &db, nil
}

type Database struct {
	Client *mongo.Client
}

func (s *Database) IsConnected() bool {
	return s.Client != nil
}

func (s *Database) ConnectWithURI(URI string) error {
	var err error
	s.Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		return err
	}
	return nil
}

func (s *Database) ConnectToCollection(database, collection string) (*mongo.Collection, error) {
	return s.Client.Database(database).Collection(collection), nil
}
