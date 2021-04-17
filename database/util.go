package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type appender func(cur *mongo.Cursor) error

func (s *Database) Aggregate(ctx context.Context, databaseName, collectionName string, condition []bson.M, appender appender) error {
	collection, err := s.ConnectToCollection(databaseName, collectionName)
	if err != nil {
		return err
	}

	options := options.Aggregate()

	cur, err := collection.Aggregate(ctx, condition, options)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		err := appender(cur)
		if err != nil {
			return err
		}
	}

	if err := cur.Err(); err != nil {
		return err
	}

	return nil
}

func (s *Database) Get(ctx context.Context, databaseName, collectionName, id string, item interface{}) error {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return s.FindOne(ctx, databaseName, collectionName, bson.M{"_id": oID}, item)
}

func (s *Database) Create(ctx context.Context, databaseName, collectionName string, item interface{}) (*primitive.ObjectID, error) {
	collection, err := s.ConnectToCollection(databaseName, collectionName)
	if err != nil {
		return nil, err
	}

	result, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	id := result.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (s *Database) Update(ctx context.Context, databaseName, collectionName string, id primitive.ObjectID, result interface{}, data interface{}) error {
	collection, err := s.ConnectToCollection(databaseName, collectionName)
	if err != nil {
		return err
	}

	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, data).Decode(result)
	return err
}

func (s *Database) Delete(ctx context.Context, databaseName, collectionName string, id string) (*mongo.DeleteResult, error) {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	collection, err := s.ConnectToCollection(databaseName, collectionName)
	if err != nil {
		return nil, err
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": oID})
	return result, err
}

func (s *Database) FindOne(ctx context.Context, databaseName, collectionName string, query bson.M, item interface{}) error {
	collection, err := s.ConnectToCollection(databaseName, collectionName)
	if err != nil {
		return err
	}

	findOptions := options.FindOne()

	return collection.FindOne(ctx, query, findOptions).Decode(item)
}

func (s *Database) Find(ctx context.Context, databaseName, collectionName string, query bson.M, limit int64, appender appender) error {
	collection, err := s.ConnectToCollection(databaseName, collectionName)
	if err != nil {
		return err
	}

	findOptions := options.Find()
	findOptions.SetLimit(limit)

	cur, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {

		err := appender(cur)
		if err != nil {
			return err
		}
	}

	if err := cur.Err(); err != nil {
		return err
	}

	return nil
}
