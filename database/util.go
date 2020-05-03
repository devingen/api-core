package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type appender func(cur *mongo.Cursor) error

func (s *Database) Aggregate(databaseName, collectionName string, condition []bson.M, appender appender) error {
	collection, err := s.ConnectToCollection(databaseName, collectionName)
	if err != nil {
		return err
	}

	options := options.Aggregate()

	cur, err := collection.Aggregate(context.TODO(), condition, options)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		err := appender(cur)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Database) Query(databaseName, collectionName string, condition bson.M, appender appender) error {
	collection, err := s.ConnectToCollection(databaseName, collectionName)
	if err != nil {
		return err
	}

	findOptions := options.Find()
	findOptions.SetLimit(10)

	cur, err := collection.Find(context.TODO(), condition, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		err := appender(cur)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Database) Find(databaseName, collectionName, id string, item interface{}) error {
	collection, err := s.ConnectToCollection(databaseName, collectionName)
	if err != nil {
		return err
	}

	findOptions := options.FindOne()

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection.FindOne(context.TODO(), bson.M{"_id": oID}, findOptions).Decode(item)
	return nil
}