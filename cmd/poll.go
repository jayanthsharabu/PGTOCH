package cmd

import (
	"context"
	"fmt"
	"pgtoch/config"
	ui "pgtoch/internal/UI"
	"pgtoch/internal/etl"
	"pgtoch/internal/log"
	"pgtoch/internal/poller"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func startPolling(ctx context.Context, cfg *config.Config, lastSeen string) error {
	log := log.StyledLog
	log.Info("Starting chg data polling..")

	startFrom := lastSeen

	if startFrom == "" {
		startFrom = "beginning"
	}

	ui.PrintBox("Polling Configuration",
		"Table: "+cfg.Table+"\n"+
			"Delta Column: "+cfg.Polling.Deltacol+"\n"+
			"Interval: "+fmt.Sprintf("%d seconds", cfg.Polling.Interval)+"\n"+
			"Starting From: "+startFrom)

	pgConn, err := pgx.Connect(ctx, cfg.PostgreSQLURL)

	if err != nil {
		fmt.Errorf("Failed to connect to PostgreSQL for polling: %w", err)
		return err
	}

	defer pgConn.Close(ctx)

	processNewData := func(data *etl.TableData) error {
		if len(data.Rows) > 0 {
			log.Info(fmt.Sprintf("Processing new batch data: %d rows", len(data.Rows)),
				zap.Int("rows", len(data.Rows)),
				zap.String("table", cfg.Table))

		} else {
			log.Info("No new data found in this cycle")
		}

		return etl.InsertRows(cfg.ClickHouseURL, cfg.Table, etl.GetColumnNames(data.Columns), data.Rows, cfg.BatchSize)
	}

	pollConfig := poller.PollConfig{
		Table:     cfg.Table,
		DeltaCol:  cfg.Polling.Deltacol,
		Interval:  time.Duration(cfg.Polling.Interval) * time.Second,
		Limit:     &cfg.Limit,
		StartFrom: lastSeen,
		OnData:    processNewData,
	}
	p := poller.NewPoller(pgConn, pollConfig)

	return p.Start(ctx)

}

func determineLastSeen(td *etl.TableData, deltaCol string) (string, error) {

	log := log.StyledLog

	if len(td.Rows) == 0 {
		log.Info("No rows")
		return "", nil
	}

	deltaColIndex := -1

	for i, col := range td.Columns {
		if col.Name == deltaCol {
			deltaColIndex = i
			break
		}
	}

	if deltaColIndex == -1 {
		return "", fmt.Errorf("delta column %s not found in table", deltaCol)
	}

	lastRow := td.Rows[len(td.Rows)-1]
	var lastSeenValue string

	switch v := lastRow[deltaColIndex].(type) {
	case string:
		lastSeenValue = v
	case int, int32, int64, int8, int16:
		lastSeenValue = fmt.Sprintf("%d", v)
	case float64, float32:
		lastSeenValue = fmt.Sprintf("%f", v)
	case time.Time:
		lastSeenValue = v.Format(time.RFC3339)
	default:
		lastSeenValue = fmt.Sprintf("%v", v)
	}

	log.Info("Determined last seen value for delta tracking",
		zap.String("column", deltaCol),
		zap.String("value", lastSeenValue))

	return lastSeenValue, nil
}
