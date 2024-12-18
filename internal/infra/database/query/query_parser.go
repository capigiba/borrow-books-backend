package query

import (
	"fmt"
	"strings"
)

func ParseFilters(rawFilters []string) ([]Filter, error) {
	var filters []Filter
	for _, f := range rawFilters {
		// expected format: field__operator__value
		parts := strings.SplitN(f, "__", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid filter format: %s", f)
		}
		field, op, val := parts[0], parts[1], parts[2]
		filters = append(filters, Filter{Field: field, Operator: op, Value: val})
	}
	return filters, nil
}

func ParseSorts(rawSorts []string) ([]Sort, error) {
	var sorts []Sort
	for _, s := range rawSorts {
		// expected format: field__asc or field__desc
		parts := strings.SplitN(s, "__", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid sort format: %s", s)
		}
		field, order := parts[0], parts[1]
		desc := false
		if order == "desc" {
			desc = true
		} else if order != "asc" {
			return nil, fmt.Errorf("invalid sort order: %s", order)
		}
		sorts = append(sorts, Sort{Field: field, Desc: desc})
	}
	return sorts, nil
}

func ParseFields(raw string) []string {
	if raw == "" {
		return nil
	}
	fields := strings.Split(raw, ",")
	for i, f := range fields {
		fields[i] = strings.TrimSpace(f)
	}
	return fields
}
