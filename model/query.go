package model

type SortConfig struct {
	ID    string `bson:"id" json:"id"`
	Order int    `bson:"order" json:"order"`
}

type BasicQueryConfig struct {
	Limit int          `bson:"limit" json:"limit"`
	Skip  int          `bson:"skip" json:"skip"`
	Sort  []SortConfig `bson:"sort" json:"sort"`
}
