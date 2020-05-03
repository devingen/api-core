package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

const (
	envMongoAddress  = "MONGO_ADDRESS"
	envMongoUser     = "MONGO_USERNAME"
	envMongoPassword = "MONGO_PASSWORD"
)

func NewDatabase() (*Database, error) {
	db := Database{}
	err := db.ConnectWithEnvironment()
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func NewDatabaseWithURI(URI string) (*Database, error) {
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
	log.Println("Connected to MongoDB")

	return nil
}

func (s *Database) ConnectWithEnvironment() error {
	address, hasVar := os.LookupEnv(envMongoAddress)
	if !hasVar {
		log.Fatalf("Missing environment variable %s", envMongoAddress)
	}
	uri := fmt.Sprintf("mongodb://%s", address)

	username := os.Getenv(envMongoUser)
	password := os.Getenv(envMongoPassword)
	if username != "" && password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s", username, password, address)
		//uri = fmt.Sprintf("mongodb+srv://%s:%s@%s", username, password, address)
	}

	var err error
	s.Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	log.Println("Connected to MongoDB")

	return nil
}

func (s *Database) ConnectToCollection(database, collection string) (*mongo.Collection, error) {
	if s.Client == nil {
		err := s.ConnectWithEnvironment()
		if err != nil {
			return nil, err
		}
	}
	return s.Client.Database(database).Collection(collection), nil
}
