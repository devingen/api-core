package model

type RelationReferenceField struct {
	Field
}

func SingleRelationReferenceFromField(field Field) RelationReferenceField {
	return RelationReferenceField{
		Field: field,
	}
}

func NewRelationReference(fieldName, relationCollection, nameInRelationCollection, otherCollection, nameOfOtherCollectionInRelationCollection string) RelationReferenceField {
	return RelationReferenceField{
		Field{
			"type":                     FieldTypeRelationReference,
			"id":                       fieldName,
			"relationCollection":       relationCollection,
			"nameInRelationCollection": nameInRelationCollection,
			"otherCollection":          otherCollection,
			"nameOfOtherCollectionInRelationCollection": nameOfOtherCollectionInRelationCollection,
		},
	}
}

func (field RelationReferenceField) GetID() string {
	return field.Field.GetID()
}

func (field RelationReferenceField) GetRelationCollection() string {
	return field.Field.GetString("relationCollection")
}

func (field RelationReferenceField) GetNameInRelationCollection() string {
	return field.Field.GetString("nameInRelationCollection")
}

func (field RelationReferenceField) GetOtherCollection() string {
	return field.Field.GetString("otherCollection")
}

func (field RelationReferenceField) GetNameOfOtherCollectionInRelationCollection() string {
	return field.Field.GetString("nameOfOtherCollectionInRelationCollection")
}

func (field RelationReferenceField) SetRelationFilter(filter *Filter) RelationReferenceField {
	field.SetInterface("relationFilter", filter)
	return field
}

func (field RelationReferenceField) GetRelationFilter() *Filter {
	return field.Field.GetFilter("relationFilter")
}

func (field RelationReferenceField) SetOtherCollectionFilter(filter *Filter) RelationReferenceField {
	field.SetInterface("otherCollectionFilter", filter)
	return field
}

func (field RelationReferenceField) GetOtherCollectionFilter() *Filter {
	return field.Field.GetFilter("otherCollectionFilter")
}

func (field RelationReferenceField) SetFields(fields []Field) RelationReferenceField {
	field.Field.SetFields(fields)
	return field
}

func (field RelationReferenceField) GetFields() []Field {
	return field.Field.GetFields()
}

func (field RelationReferenceField) SetOtherCollectionFields(fields []Field) RelationReferenceField {
	field.Field.SetFieldsForKey("otherCollectionFields", fields)
	return field
}

func (field RelationReferenceField) GetOtherCollectionFields() []Field {
	return field.Field.GetFieldsForKey("otherCollectionFields")
}

func (field RelationReferenceField) ToField() Field {
	return field.Field
}
