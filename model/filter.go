package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Operator string

const (
	OperatorAnd Operator = "and"
	OperatorOr  Operator = "or"
)

type Comparison string

const (
	ComparisonEqMongoOID Comparison = "eq-mongo-oid"
	ComparisonEq         Comparison = "eq"
	ComparisonNe         Comparison = "ne"
	ComparisonLt         Comparison = "lt"
	ComparisonLte        Comparison = "lte"
	ComparisonGt         Comparison = "gt"
	ComparisonGte        Comparison = "gte"
	ComparisonIn         Comparison = "in"
	ComparisonContain    Comparison = "contain"
	ComparisonNcontain   Comparison = "ncontain"
	ComparisonSimilar    Comparison = "similar"
	ComparisonDifferent  Comparison = "different"
	ComparisonEmpty      Comparison = "empty"
	ComparisonNotEmpty   Comparison = "nempty"
)

type Filter struct {
	Comparison Comparison  `json:"comparison"`
	FieldId    string      `json:"id"`
	FieldValue interface{} `json:"value"`

	Filters  []Filter `json:"filters"`
	Operator Operator `json:"operator"`
}

func (c Filter) ToMatchQuery(config *QueryConfig) bson.M {
	if c.Filters != nil {
		conditions := make([]bson.M, len(c.Filters))
		for i, condition := range c.Filters {
			conditions[i] = condition.ToMatchQuery(config)
		}
		return bson.M{"$" + string(c.Operator): conditions}
	}

	if c.Comparison == ComparisonNotEmpty {
		return bson.M{c.FieldId: bson.M{"$exists": true}}
	} else if c.Comparison == ComparisonEmpty {
		return bson.M{c.FieldId: bson.M{"$exists": false}}
	}

	field := config.GetField(c.FieldId)
	if field == nil {
		// it may be a filter for an inner field of a relation. Ex: { "fieldId": "organisation._id" }
		if c.Comparison == ComparisonContain {
			return bson.M{c.FieldId: bson.M{"$regex": c.FieldValue, "$options": "i"}}
		}
		if c.Comparison == ComparisonEqMongoOID {
			oid, _ := primitive.ObjectIDFromHex(c.FieldValue.(string))
			return bson.M{c.FieldId: bson.M{"$" + string(ComparisonEq): oid}}
		}
		return bson.M{c.FieldId: bson.M{"$" + string(c.Comparison): c.FieldValue}}
	}
	if field.GetType() == FieldTypeText {
		if c.Comparison == ComparisonContain {
			return bson.M{c.FieldId: bson.M{"$regex": c.FieldValue, "$options": "i"}}
		}
		if c.Comparison == ComparisonNcontain {
			return bson.M{c.FieldId: bson.M{"$regex": "^((?!" + c.FieldValue.(string) + ").)*$", "$options": "i"}}
		}
		return bson.M{c.FieldId: bson.M{"$" + string(c.Comparison): c.FieldValue}}
	}
	if field.GetType() == FieldTypeNumber {
		return bson.M{c.FieldId: bson.M{"$" + string(c.Comparison): c.FieldValue}}
	}
	return bson.M{c.FieldId: bson.M{"$regex": c.FieldValue, "$options": "i"}}
}

func (c Filter) ToFilterQuery(name string) bson.M {
	if c.Filters != nil {
		conditions := make([]bson.M, len(c.Filters))
		for i, condition := range c.Filters {
			conditions[i] = condition.ToFilterQuery(name)
		}
		return bson.M{"$" + string(c.Operator): conditions}
	}
	return bson.M{"$" + string(c.Comparison): []interface{}{"$$" + name + "." + c.FieldId, c.FieldValue}}
}

func FilterFromMap(data map[string]interface{}) *Filter {
	filtersMap, hasFilters := data["filters"].([]interface{})
	if hasFilters {
		filters := make([]Filter, len(filtersMap))
		for i, filterMap := range filtersMap {
			filters[i] = *FilterFromMap(filterMap.(map[string]interface{}))
		}
		return &Filter{
			Filters:  filters,
			Operator: Operator(data["operator"].(string)),
		}
	}
	return &Filter{
		Comparison: Comparison(data["comparison"].(string)),
		FieldId:    data["id"].(string),
		FieldValue: data["value"].(string),
	}
}
