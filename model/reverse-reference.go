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
