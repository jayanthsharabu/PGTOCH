package etl

import (
	"context"
	"fmt"
	"pgtoch/internal/db"
	"pgtoch/internal/log"
	"time"

	"go.uber.org/zap"
)

func CreateTable(chURL, ddl string) error {
	conn, err := db.ConnectClickhouse(chURL)
	if err != nil {
		return fmt.Errorf("failed to connect to clickhouse: %w", err)
	}
	defer conn.Close()

	_, err = conn.ExecContext(context.Background(), ddl)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	log.StyledLog.Success("Table created successfully")
	return nil
}

func InsertRows(chURL, table string, columns []string, rows [][]any, batchSize int) error {
	if !IsValidIdentifier(table) {
		return fmt.Errorf("invalid table name: %s", table)
	}

	for _, col := range columns {
		if !IsValidIdentifier(col) {
			return fmt.Errorf("invalid column name: %s", col)
		}
	}

	conn, err := db.ConnectClickhouse(chURL)
	if err != nil {
		return fmt.Errorf("failed to connect to clickhouse: %w", err)
	}
	defer conn.Close()

	quotedColumns := make([]string, len(columns))
	for i, col := range columns {
		quotedColumns[i] = QuoteIdentifier(col)
	}

	colNames := "(" + join(quotedColumns, ", ") + ")"

	quotedTable := QuoteIdentifier(table)

	insertPrefix := fmt.Sprintf("INSERT INTO %s %s VALUES", quotedTable, colNames)

	ctx := context.Background()

	for i := 0; i < len(rows); i += batchSize {
		end := min(i+batchSize, len(rows))
		batch := rows[i:end]

		query := insertPrefix + buildValuesPlaceholders(len(batch), len(columns))
		args := flatten(batch)

		err := Retry(ctx, RetryConfig{
			MaxAttempts: 4,
			BaseDelay:   250 * time.Millisecond,
			MaxDelay:    10 * time.Second,
			Jitter:      true,
		}, func() error {
			_, err := conn.ExecContext(ctx, query, args...)
			if err != nil {
				return fmt.Errorf("failed to insert batch: %w", err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to insert rows into %s: %w", table, err)
		}

		log.Logger.Info("Inserted Into Clickhouse",
			zap.Int("row_count", end-i),
			zap.String("table", table),
			zap.Int("batch_size", batchSize),
			zap.Int("total_rows", len(rows)),
		)

	}
	return nil

}
