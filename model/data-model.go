package model

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//type DataModel struct {
//	ID        primitive.ObjectID `json:"_id" bson:"_id"`
//	CreatedAt time.Time          `json:"_created" bson:"_created"`
//	UpdatedAt time.Time          `json:"_updated" bson:"_updated"`
//	Revision  int                `json:"_revision" bson:"_revision,omitempty"`
//	Domain    string             `json:"domain" bson:"domain"`
//}

type DataModel map[string]interface{}

func (dm DataModel) GetID() string {
	id, has := dm["_id"].(string)
	if !has {
		return ""
	}
	return id
}

func (dm DataModel) GetFieldCount() int {
	return len(dm)
}

func (dm DataModel) GetDBRef(key string) *DBRef {
	dbref, has := dm[key].(DBRef)
	if !has {
		dbrefAsMap, hasRefAsMap := dm[key].(map[string]interface{})
		if !hasRefAsMap {
			return nil
		}
		id, _ := primitive.ObjectIDFromHex(dbrefAsMap["$id"].(string))
		return &DBRef{
			Ref:      dbrefAsMap["$ref"].(string),
			ID:       id,
			Database: dbrefAsMap["$db"].(string),
		}
	}
	return &dbref
}

func (dm DataModel) GetChildDataModel(key string) DataModel {
	model, has := dm[key].(DataModel)
	if !has {
		model, has = dm[key].(map[string]interface{})
		if !has {
			return nil
		}
		return model
	}
	return model
}

func (dm DataModel) GetChildDataModelArray(key string) []DataModel {
	models, has := dm[key].([]DataModel)
	if !has {
		modelsAsMap, hasModelsAsMap := dm[key].(primitive.A)
		if !hasModelsAsMap {
			return nil
		}
		models := make([]DataModel, len(modelsAsMap))
		for i, m := range modelsAsMap {
			models[i] = m.(DataModel)
		}
		return models
	}
	return models
}

func (dm DataModel) GetString(key string) string {
	value, has := dm[key].(string)
	if !has {
		return ""
	}
	return value
}

func (dm DataModel) GetFields() []Field {
	return dm.GetFieldsForKey("fields")
}

func (dm DataModel) SetFields(fields []Field) DataModel {
	dm.SetFieldsForKey("fields", fields)
	return dm
}

func (dm DataModel) SetFieldsForKey(key string, fields []Field) DataModel {
	dm[key] = fields
	return dm
}

func (dm DataModel) GetFieldsForKey(key string) []Field {
	value, has := dm.GetInterface(key).([]Field)
	if !has {
		valueInterface, hasInterface := dm.GetInterface(key).([]interface{})
		if !hasInterface {
			return nil
		}
		fields := make([]Field, len(valueInterface))
		for i, filterMap := range valueInterface {
			fields[i] = filterMap.(map[string]interface{})
		}
		return fields
	}
	return value
}

func (dm DataModel) GetInterface(key string) interface{} {
	return dm[key]
}

// ToStruct converts DataModel into T.
// Ex: core_model.ToStruct[model.UserMeta](response.Results[0])
func ToStruct[T any](dm *DataModel) *T {
	bytes, err := json.Marshal(dm)
	if err != nil {
		return nil
	}

	var result T
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil
	}
	return &result
}

// ToStructList converts DataModel list into T list.
// Ex: core_model.ToStructList[model.User](response.Results)
func ToStructList[T any](dmList []*DataModel) []T {
	list := make([]T, len(dmList))
	for i, dm := range dmList {
		list[i] = *ToStruct[T](dm)
	}
	return list
}
