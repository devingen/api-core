package model

type ReferenceField struct {
	Field
}

func ReferenceFromField(field Field) ReferenceField {
	return ReferenceField{
		Field: field,
	}
}

func NewReference(fieldName, otherCollection string, isSingle bool) ReferenceField {
	return ReferenceField{
		Field: Field{
			"type":            FieldTypeReference,
			"id":              fieldName,
			"otherCollection": otherCollection,
			"isSingle":        isSingle,
		},
	}
}

func (field ReferenceField) GetID() string {
	return field.Field.GetID()
}

func (field ReferenceField) GetOtherCollection() string {
	return field.Field.GetString("otherCollection")
}

func (field ReferenceField) IsSingle() bool {
	return field.Field.GetBool("isSingle")
}

func (field ReferenceField) SetFilter(filter *Filter) ReferenceField {
	field.SetInterface("filter", filter)
	return field
}

func (field ReferenceField) GetFilter() *Filter {
	return field.Field.GetFilter("filter")
}

func (field ReferenceField) SetFields(fields []Field) ReferenceField {
	field.Field.SetFields(fields)
	return field
}

func (field ReferenceField) GetFields() []Field {
	return field.Field.GetFields()
}

func (field ReferenceField) ToField() Field {
	return field.Field
}
