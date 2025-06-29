package etl

import (
	"fmt"
	"strings"
)

func MapColumnType(cols []Column) ([]string, error) {
	var mapped []string
	for _, col := range cols {
		chType, ok := pgtochtype[col.Type]
		if !ok {
			return nil, fmt.Errorf("unsupported column type: %s", col.Type)
		}
		mapped = append(mapped, fmt.Sprintf("%s %s", col.Name, chType))
	}
	return mapped, nil
}

func BuildDDLQuery(table string, cols []Column) (string, error) {
	mappedCols, err := MapColumnType(cols)
	if err != nil {
		return "", err
	}
	ddl := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s) ENGINE = MergeTree() ORDER BY tuple();", QuoteIdentifier(table), strings.Join(mappedCols, ", "))
	if len(mappedCols) == 0 {
		return "", fmt.Errorf("no columns to create table")
	}
	return ddl, nil
}
