package model

type DBRef struct {
	Ref      string `json:"$ref" bson:"$ref"`
	ID       string `json:"$id" bson:"$id"`
	Database string `json:"$db" bson:"$db"`
}
