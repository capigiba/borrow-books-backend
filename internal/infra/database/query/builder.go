package query

import (
	"fmt"
	"strings"
)

func BuildSelectQuery(tableName string, opts QueryOptions) (string, []interface{}) {
	var (
		args         []interface{}
		whereClauses []string
		selectFields string
		orderClause  string
	)

	// Handle fields
	if len(opts.Fields) == 0 {
		selectFields = "*"
	} else {
		selectFields = strings.Join(opts.Fields, ", ")
	}

	// Handle filters
	argIndex := 1
	for _, fil := range opts.Filters {
		op := ""
		switch fil.Operator {
		case "eq":
			op = "="
		case "ilike":
			op = "ILIKE"
		case "gt":
			op = ">"
		case "gte":
			op = ">="
		case "lt":
			op = "<"
		case "lte":
			op = "<="
		default:
			op = "="
		}

		whereClauses = append(whereClauses, fmt.Sprintf("%s %s $%d", fil.Field, op, argIndex))
		args = append(args, fil.Value)
		argIndex++
	}

	// Handle sorts
	if len(opts.Sorts) > 0 {
		var sortExprs []string
		for _, s := range opts.Sorts {
			dir := "ASC"
			if s.Desc {
				dir = "DESC"
			}
			sortExprs = append(sortExprs, fmt.Sprintf("%s %s", s.Field, dir))
		}
		orderClause = "ORDER BY " + strings.Join(sortExprs, ", ")
	}

	query := fmt.Sprintf("SELECT %s FROM %s", selectFields, tableName)
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	if orderClause != "" {
		query += " " + orderClause
	}
	return query, args
}
