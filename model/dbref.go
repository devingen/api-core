package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type DBRef struct {
	Ref      string             `bson:"_ref" json:"_ref"`
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Database string             `bson:"_db" json:"_db"`
}
