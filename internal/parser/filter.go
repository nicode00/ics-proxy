package parser

import (
	"strings"
)

type Filter struct {
	Field    string
	Operator string
	Arg      string
}

func (f Filter) FieldsMatch(fields map[string]string) bool {
	val, exists := fields[f.Field]

	switch f.Operator {
	case "less_than":
		return exists && val < f.Arg
	case "greater_than":
		return exists && val > f.Arg
	case "equals":
		return exists && val == f.Arg
	case "not_equals":
		return !exists || val != f.Arg
	case "contains":
		return exists && strings.Contains(val, f.Arg)
	case "not_contains":
		return !exists || !strings.Contains(val, f.Arg)
	default:
		panic("unrecognized operator `" + f.Operator + "`")
	}
}
