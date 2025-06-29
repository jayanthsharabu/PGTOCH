package export

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"pgtoch/internal/db"
	"pgtoch/internal/etl"
	"pgtoch/internal/log"

	"go.uber.org/zap"
)

func ExportTabletoCSV(chURL, table, outPath string) error {
	conn, err := db.ConnectClickhouse(chURL)
	if err != nil {
		return fmt.Errorf("failed to connect to clickhouse: %w", err)
	}
	defer conn.Close()

	exists, err := TableExists(conn, table)
	if err != nil {
		return fmt.Errorf("failed to check if table exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("table %s does not exist", table)
	}
	ctx := context.Background()

	rows, err := conn.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %s", etl.QuoteIdentifier(table)))
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}
	if len(cols) == 0 {
		return fmt.Errorf("no columns found in table %s", table)
	}

	file, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(cols); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	count := 0

	for rows.Next() {
		columns := make([]any, len(cols))
		columnPtrs := make([]any, len(cols))

		for i := range columns {
			columnPtrs[i] = &columns[i]
		}
		if err := rows.Scan(columnPtrs...); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		record := make([]string, len(cols))
		for i, col := range columns {
			if col == nil {
				record[i] = ""
			} else {
				record[i] = fmt.Sprintf("%v", col)
			}
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("could not write to csv %w", err)
		}
		count++
		if count%500 == 0 {
			log.Logger.Info("Exported rows ..\n",
				zap.Int("rows_count", count),
				zap.String("table", table),
			)
		}
	}
	log.Logger.Info("Exported rows ..\n",
		zap.String("table", table),
		zap.Int("rows_count", count),
		zap.String("outPath", outPath),
	)
	return nil
}
