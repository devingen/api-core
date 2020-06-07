package model

type ReverseReferenceField struct {
	Field
}

func ReverseReferenceFromField(field Field) ReverseReferenceField {
	return ReverseReferenceField{
		Field: field,
	}
}

func NewReverseReference(fieldName, otherCollection, nameInOtherCollection string, isSingle bool) ReverseReferenceField {
	return ReverseReferenceField{
		Field: Field{
			"type":                  FieldTypeReverseReference,
			"id":                    fieldName,
			"otherCollection":       otherCollection,
			"nameInOtherCollection": nameInOtherCollection,
			"isSingle":              isSingle,
		},
	}
}

func (field ReverseReferenceField) GetID() string {
	return field.Field.GetID()
}

func (field ReverseReferenceField) GetOtherCollection() string {
	return field.Field.GetString("otherCollection")
}

func (field ReverseReferenceField) GetNameInOtherCollection() string {
	return field.Field.GetString("nameInOtherCollection")
}

func (field ReverseReferenceField) IsSingle() bool {
	return field.Field.GetBool("isSingle")
}

func (field ReverseReferenceField) SetFilter(filter *Filter) ReverseReferenceField {
	field.SetInterface("filter", filter)
	return field
}

func (field ReverseReferenceField) GetFilter() *Filter {
	return field.Field.GetFilter("filter")
}

func (field ReverseReferenceField) SetLimit(limit int) ReverseReferenceField {
	field.Field.SetInt("limit", limit)
	return field
}

func (field ReverseReferenceField) GetLimit() int {
	return field.Field.GetInt("limit")
}

func (field ReverseReferenceField) GetSort() *SortConfig {
	value, has := field.GetInterface("sort").(*SortConfig)
	if !has {
		valueInterface, hasInterface := field.GetInterface("sort").(map[string]interface{})
		if !hasInterface {
			return nil
		}
		return &SortConfig{
			ID: valueInterface["id"].(string),
			Order: int(valueInterface["order"].(float64)),
		}
	}
	return value
}

func (field ReverseReferenceField) SetSort(sort *SortConfig) ReverseReferenceField {
	field.Field["sort"] = sort
	return field
}

func (field ReverseReferenceField) SetFields(fields []Field) ReverseReferenceField {
	field.Field.SetFields(fields)
	return field
}

func (field ReverseReferenceField) GetFields() []Field {
	return field.Field.GetFields()
}

func (field ReverseReferenceField) ToField() Field {
	return field.Field
}
