package etl

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Column struct {
	Name string
	Type string
}

type TableData struct {
	Columns []Column
	Rows    [][]any
}

func getColumns(ctx context.Context, conn *pgx.Conn, table string) ([]Column, error) {
	colQuery := `
	SELECT column_name, data_type
	FROM information_schema.columns
	WHERE table_name = $1
	ORDER BY ordinal_position
	`
	rows, err := conn.Query(ctx, colQuery, table)

	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()
	var cols []Column
	for rows.Next() {
		var col Column
		if err := rows.Scan(&col.Name, &col.Type); err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}
		cols = append(cols, col)
	}
	return cols, nil

}

func ExtractTableData(ctx context.Context, conn *pgx.Conn, table string, limit *int) (*TableData, error) {

	cols, err := getColumns(ctx, conn, table)
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var query string
	var rows pgx.Rows

	if limit != nil && *limit > 0 {
		query = "SELECT * FROM " + pgx.Identifier{table}.Sanitize() + "LIMIT $1"
		rows, err = conn.Query(ctx, query, *limit)
	} else {
		query = "SELECT * FROM " + pgx.Identifier{table}.Sanitize()
		rows, err = conn.Query(ctx, query)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	defer rows.Close()

	var results [][]any

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, fmt.Errorf("failed to get row values: %w", err)
		}
		for i, val := range values {
			var uuidBytes []byte
			switch v := val.(type) {
			case [16]byte:
				uuidBytes = v[:]
			case []byte:
				uuidBytes = v
			}

			if uuidBytes != nil && (cols[i].Type == "uuid" || cols[i].Type == "bytea") {
				if len(uuidBytes) == 16 {
					values[i] = fmt.Sprintf("%x-%x-%x-%x-%x", uuidBytes[0:4], uuidBytes[4:6], uuidBytes[6:8], uuidBytes[8:10], uuidBytes[10:16])
				}
			}
		}
		results = append(results, values)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	return &TableData{
		Columns: cols,
		Rows:    results,
	}, nil
}

func ExtractTableDataSince(ctx context.Context, conn *pgx.Conn, table, deltaCol, lastSeen string, limit *int) (*TableData, error) {

	cols, err := getColumns(ctx, conn, table)
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE %s > $1 ORDER BY %s ASC",
		pgx.Identifier{table}.Sanitize(),
		deltaCol,
		deltaCol,
	)
	var rows pgx.Rows
	if limit != nil && *limit > 0 {
		query += fmt.Sprintf("LIMIT $1")
		rows, err = conn.Query(ctx, query, lastSeen, *limit)
	} else {
		rows, err = conn.Query(ctx, query, lastSeen)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	defer rows.Close()

	var results [][]any

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, fmt.Errorf("failed to get row values: %w", err)
		}
		for i, val := range values {
			var uuidBytes []byte
			switch v := val.(type) {
			case [16]byte:
				uuidBytes = v[:]
			case []byte:
				uuidBytes = v
			}

			if uuidBytes != nil && (cols[i].Type == "uuid" || cols[i].Type == "bytea") {
				if len(uuidBytes) == 16 {
					values[i] = fmt.Sprintf("%x-%x-%x-%x-%x", uuidBytes[0:4], uuidBytes[4:6], uuidBytes[6:8], uuidBytes[8:10], uuidBytes[10:16])
				}
			}
		}
		results = append(results, values)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating delta rows: %w", err)
	}

	return &TableData{
		Columns: cols,
		Rows:    results,
	}, nil

}

func GetColumnNames(cols []Column) []string {
	names := make([]string, len(cols))
	for i, col := range cols {
		names[i] = col.Name
	}
	return names
}
