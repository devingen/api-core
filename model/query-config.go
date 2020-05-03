package model

type QueryConfig struct {
	Filter *Filter      `bson:"filter" json:"filter"`
	Limit  int          `bson:"limit" json:"limit"`
	Skip   int          `bson:"skip" json:"skip"`
	Sort   []SortConfig `bson:"sort" json:"sort"`
	Fields []Field      `bson:"fields" json:"fields"`
}

func (c *QueryConfig) GetField(id string) *Field {
	for _, field := range c.Fields {
		if field.GetID() == id {
			return &field
		}
	}
	return nil
}
