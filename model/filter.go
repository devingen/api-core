package model

import (
	"strconv"
	"time"

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

	ComparisonDateEqDay     Comparison = "date-eq-day"
	ComparisonDateNeDay     Comparison = "date-ne-day"
	ComparisonDateNextYear  Comparison = "date-next-year"
	ComparisonDateNextMonth Comparison = "date-next-month"
	ComparisonDateNextWeek  Comparison = "date-next-week"
	ComparisonDateThisYear  Comparison = "date-this-year"
	ComparisonDateThisMonth Comparison = "date-this-month"
	ComparisonDateThisWeek  Comparison = "date-this-week"
	ComparisonDateLastYear  Comparison = "date-last-year"
	ComparisonDateLastMonth Comparison = "date-last-month"
	ComparisonDateLastWeek  Comparison = "date-last-week"
	DateNextNumberOfDays    Comparison = "date-next-number-of-days"
	DatePastNumberOfDays    Comparison = "date-past-number-of-days"
)

type Filter struct {
	Comparison Comparison  `json:"comparison"`
	FieldId    string      `json:"id"`
	FieldValue interface{} `json:"value"`

	Filters  []Filter `json:"filters"`
	Operator Operator `json:"operator"`
}

var specialDateConditions = map[Comparison]bool{
	ComparisonDateEqDay:     true,
	ComparisonDateNeDay:     true,
	DateNextNumberOfDays:    true,
	DatePastNumberOfDays:    true,
	ComparisonDateNextYear:  true,
	ComparisonDateNextMonth: true,
	ComparisonDateNextWeek:  true,
	ComparisonDateThisYear:  true,
	ComparisonDateThisMonth: true,
	ComparisonDateThisWeek:  true,
	ComparisonDateLastYear:  true,
	ComparisonDateLastMonth: true,
	ComparisonDateLastWeek:  true,
}

// nonNegativeIntFromValue converts various JSON-decoded numeric representations
// (int, float64, string) into a non-negative int. If parsing fails, returns 0.
func nonNegativeIntFromValue(v interface{}) int {
	var n int
	switch t := v.(type) {
	case int:
		n = t
	case float64:
		// numbers parsed from JSON may come as float64
		n = int(t)
	case string:
		if parsed, err := strconv.Atoi(t); err == nil {
			n = parsed
		}
	}
	if n < 0 {
		n = 0
	}
	return n
}

func getDateFilters(filter Filter) (start time.Time, end time.Time) {
	base := time.Now()
	loc := base.Location()
	switch filter.Comparison {
	case ComparisonDateEqDay:
		t, err := time.Parse(time.RFC3339, filter.FieldValue.(string))
		if err == nil {
			start = time.Date(base.Year(), base.Month(), t.Day(), 0, 0, 0, 0, loc)
			end = time.Date(base.Year(), base.Month(), t.Day(), 23, 59, 59, 999999999, loc)
		}
	case DateNextNumberOfDays:
		// Compute range from start of today to end of the day after adding N days (inclusive)
		daysInt := nonNegativeIntFromValue(filter.FieldValue)
		startOfToday := time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, loc)
		start = startOfToday
		end = startOfToday.AddDate(0, 0, daysInt+1).Add(-time.Nanosecond)
	case DatePastNumberOfDays:
		// Compute range from start of the day N days ago to end of today (inclusive)
		daysInt := nonNegativeIntFromValue(filter.FieldValue)
		startOfToday := time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, loc)
		start = startOfToday.AddDate(0, 0, -daysInt)
		end = startOfToday.AddDate(0, 0, 1).Add(-time.Nanosecond)
	case ComparisonDateNextYear:
		y := base.Year() + 1
		start = time.Date(y, time.January, 1, 0, 0, 0, 0, loc)
		end = time.Date(y, time.December, 31, 23, 59, 59, 999999999, loc)
	case ComparisonDateNextMonth:
		y, m := base.Year(), base.Month()
		if m == time.December {
			y++
			m = time.January
		} else {
			m++
		}
		start = time.Date(y, m, 1, 0, 0, 0, 0, loc)
		// end is the last moment of that next month
		y2, m2 := y, m
		if m2 == time.December {
			y2++
			m2 = time.January
		} else {
			m2++
		}
		end = time.Date(y2, m2, 1, 0, 0, 0, 0, loc).Add(-time.Nanosecond)
	case ComparisonDateNextWeek:
		b := base.AddDate(0, 0, 7)
		weekday := int(b.Weekday())
		daysSinceMonday := (weekday + 6) % 7
		weekStart := time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -daysSinceMonday)
		start = weekStart
		end = weekStart.AddDate(0, 0, 7).Add(-time.Nanosecond)
	case ComparisonDateThisYear:
		start = time.Date(base.Year(), time.January, 1, 0, 0, 0, 0, loc)
		end = time.Date(base.Year(), time.December, 31, 23, 59, 59, 999999999, loc)
	case ComparisonDateThisMonth:
		y, m := base.Year(), base.Month()
		start = time.Date(y, m, 1, 0, 0, 0, 0, loc)
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	case ComparisonDateThisWeek:
		weekday := int(base.Weekday())
		daysSinceMonday := (weekday + 6) % 7
		weekStart := time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -daysSinceMonday)
		start = weekStart
		end = weekStart.AddDate(0, 0, 7).Add(-time.Nanosecond)
	case ComparisonDateLastYear:
		start = time.Date(base.Year()-1, time.January, 1, 0, 0, 0, 0, loc)
		end = time.Date(base.Year()-1, time.December, 31, 23, 59, 59, 999999999, loc)
	case ComparisonDateLastMonth:
		y, m := base.Year(), base.Month()
		if m == time.January {
			y--
			m = time.December
		} else {
			m--
		}
		start = time.Date(y, m, 1, 0, 0, 0, 0, loc)
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	case ComparisonDateLastWeek:
		b := base.AddDate(0, 0, -7)
		weekday := int(b.Weekday())
		daysSinceMonday := (weekday + 6) % 7
		weekStart := time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -daysSinceMonday)
		start = weekStart
		end = weekStart.AddDate(0, 0, 7).Add(-time.Nanosecond)
	}
	return
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
		numberValue := nonNegativeIntFromValue(c.FieldValue)
		return bson.M{c.FieldId: bson.M{"$" + string(c.Comparison): numberValue}}
	}
	if field.GetType() == FieldTypeBoolean {
		return bson.M{c.FieldId: bson.M{"$" + string(c.Comparison): c.FieldValue}}
	}
	if field.GetType() == FieldTypeDate {
		if c.Comparison == ComparisonDateNeDay {
			t, err := time.Parse(time.RFC3339, c.FieldValue.(string))
			if err == nil {
				base := time.Now()
				loc := base.Location()
				start := time.Date(base.Year(), base.Month(), t.Day(), 0, 0, 0, 0, loc)
				end := time.Date(base.Year(), base.Month(), t.Day(), 23, 59, 59, 999999999, loc)
				return bson.M{
					"$or": []bson.M{
						{c.FieldId: bson.M{"$lt": start}}, // before start
						{c.FieldId: bson.M{"$gt": end}},   // after end
					},
				}
			}
		}
		if specialDateConditions[c.Comparison] {
			start, end := getDateFilters(c)
			return bson.M{c.FieldId: bson.M{"$gte": start, "$lte": end}}
		}

		// Fallback to standard single-date comparisons (eq, lt, lte, gt, gte, empty, nempty)
		t, err := time.Parse(time.RFC3339, c.FieldValue.(string))
		if err == nil {
			return bson.M{c.FieldId: bson.M{"$" + string(c.Comparison): t}}
		}
	}
	if c.Comparison == ComparisonNe {
		return bson.M{c.FieldId: bson.M{"$ne": c.FieldValue}}
	}

	// Who uses this default step?
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
