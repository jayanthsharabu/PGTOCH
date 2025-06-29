package export

import (
	"context"
	"database/sql"
	"fmt"
	"pgtoch/internal/etl"
)

func TableExists(conn *sql.DB, table string) (bool, error) {
	ctx := context.Background()
	var exists uint8
	err := conn.QueryRowContext(ctx, fmt.Sprintf("EXISTS TABLE %s", etl.QuoteIdentifier(table))).Scan(&exists)
	return exists == 1, err
}
