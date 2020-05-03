package model

type CollectionLookupField struct {
	Field
}

func CollectionLookupFieldFromField(field Field) CollectionLookupField {
	return CollectionLookupField{
		Field: field,
	}
}

func NewCollectionLookup(from, foreignField, localField string, isSingle bool) CollectionLookupField {
	return CollectionLookupField{
		Field: Field{
			"type":         FieldTypeCollectionLookup,
			"from":         from,
			"foreignField": foreignField,
			"localField":   localField,
			"isSingle":     isSingle,
		},
	}
}

func (field CollectionLookupField) GetFrom() string {
	return field.Field.GetString("from")
}

func (field CollectionLookupField) GetLocalField() string {
	return field.Field.GetString("localField")
}

func (field CollectionLookupField) GetForeignField() string {
	return field.Field.GetString("foreignField")
}

func (field CollectionLookupField) IsSingle() bool {
	return field.Field.GetBool("isSingle")
}

func (field CollectionLookupField) ToField() Field {
	return field.Field
}
