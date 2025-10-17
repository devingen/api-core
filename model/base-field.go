package model

type Field map[string]interface{}

type FieldType string

const (
	FieldTypeAny               FieldType = "any" // used to fetch data without manipulating the inner data or type
	FieldTypeBoolean           FieldType = "boolean"
	FieldTypeText              FieldType = "text"
	FieldTypeNumber            FieldType = "number"
	FieldTypeReference         FieldType = "reference"
	FieldTypeReverseReference  FieldType = "reverse-reference"
	FieldTypeRelationReference FieldType = "relation-reference"
	FieldTypeCollectionLookup  FieldType = "collection-lookup"
	FieldTypeDataTransfer      FieldType = "data-transfer"
)

func New(fieldType FieldType, id string) Field {
	return Field{
		"type": fieldType,
		"id":   id,
	}
}

func (field Field) GetID() string {
	id, has := field["id"].(string)
	if !has {
		return ""
	}
	return id
}

func (field Field) GetType() FieldType {
	fieldType, has := field["type"].(FieldType)
	if !has {
		fieldTypeString, hasString := field["type"].(string)
		if !hasString {
			return ""
		}
		return FieldType(fieldTypeString)
	}
	return fieldType
}

func (field Field) GetFilter(name string) *Filter {
	value, has := field.GetInterface(name).(*Filter)
	if !has {
		valueInterface, hasInterface := field.GetInterface(name).(map[string]interface{})
		if !hasInterface {
			return nil
		}
		return FilterFromMap(valueInterface)
	}
	return value
}

func (field Field) SetString(key, value string) Field {
	field[key] = value
	return field
}

func (field Field) GetString(key string) string {
	value, has := field[key].(string)
	if !has {
		return ""
	}
	return value
}

func (field Field) SetInt(key string, value int) Field {
	field[key] = value
	return field
}

func (field Field) GetInt(key string) int {

	switch v := field[key].(type) {
	case int:
		return v
	case float64:
		// int values are populated as float64 when it's parsed from JSON
		return int(v)
	default:
		return 0
	}
}

func (field Field) GetBool(key string) bool {
	value, has := field[key].(bool)
	if !has {
		return false
	}
	return value
}

func (field Field) SetInterface(key string, value interface{}) Field {
	field[key] = value
	return field
}

func (field Field) GetInterface(key string) interface{} {
	return field[key]
}

func (field Field) SetFields(fields []Field) Field {
	field.SetFieldsForKey("fields", fields)
	return field
}

func (field Field) GetFields() []Field {
	return field.GetFieldsForKey("fields")
}

func (field Field) SetFieldsForKey(key string, fields []Field) Field {
	field[key] = fields
	return field
}

func (field Field) GetFieldsForKey(key string) []Field {
	value, has := field.GetInterface(key).([]Field)
	if !has {
		valueInterface, hasInterface := field.GetInterface(key).([]interface{})
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
