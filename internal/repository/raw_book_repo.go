package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ExtraRepository interface {
	ExecuteRawQuery(ctx context.Context, query string) ([]map[string]interface{}, error)
}

type extraRepository struct {
	db *sqlx.DB
}

func NewExtraRepository(db *sqlx.DB) ExtraRepository {
	return &extraRepository{db: db}
}

func (r *extraRepository) ExecuteRawQuery(ctx context.Context, query string) ([]map[string]interface{}, error) {
	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		// Scan each row into a map
		rowMap, err := rows.SliceScan()
		if err != nil {
			return nil, err
		}

		// Convert row slice into a map with column names as keys
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		rowData := make(map[string]interface{})
		for i, col := range columns {
			rowData[col] = rowMap[i]
		}
		results = append(results, rowData)
	}

	return results, rows.Err()
}
