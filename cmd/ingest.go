package cmd

import (
	"context"
	"fmt"
	"pgtoch/config"
	ui "pgtoch/internal/UI"
	"pgtoch/internal/db"
	"pgtoch/internal/etl"
	"pgtoch/internal/log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	ingestPgURL, ingestChURL, ingestTable, ingestConfigPath, ingestPollDelta string
	ingestLimit, ingestBatch, ingestPollInt                                  int
	ingestPoll                                                               bool
)

var ingestCmd = &cobra.Command{
	Use:   "data ingest",
	Short: "transeferring postgres to clickhouse",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintTitle("Data Ingestion")
		ui.PrintSubtitle("transferring postgres to clickhouse")

		ctx := context.Background()
		log := log.StyledLog
		log.Info("Starting data ingestion..")

		cfg := loadConfig()

		if !validateConfig(cfg) {
			return
		}

		ui.PrintBox("Configuration",
			"PostgreSQL: Connected\n"+
				"ClickHouse: Connected\n"+
				"Target Table: "+cfg.Table+"\n"+
				"Batch Size: "+ui.HighlightStyle.Render(UI_itoa(cfg.BatchSize))+" rows\n"+
				"Limit: "+ui.HighlightStyle.Render(UI_itoa(cfg.Limit))+" rows")

		conn, err := db.ConnectPostgres(cfg.PostgreSQLURL)
		if err != nil {
			log.Error("Failed to connect to PostgreSQL", zap.Error(err))
			return
		}
		defer conn.Close(ctx)

		log.Info("extracting table data")

		td, err := etl.ExtractTableData(ctx, conn, cfg.Table, &cfg.Limit)

		if err != nil {
			log.Error("failed to extract data from table", zap.Error(err))
			return
		}

		log.Success("Extracted Table data",
			zap.String("table", cfg.Table),
			zap.Int("rows", len(td.Rows)),
			zap.Int("columns", len(td.Columns)),
		)

		log.Info("Building ClickHouse schema")

		ddl, err := etl.BuildDDLQuery(cfg.Table, td.Columns)
		if err != nil {
			log.Error("failed to build DDL query", zap.Error(err))
			return
		}

		chConn, err := db.ConnectClickhouse(cfg.ClickHouseURL)
		if err != nil {
			log.Error("failed to connect to ClickHouse", zap.Error(err))
			return
		}
		defer chConn.Close()

		log.Info("creating table in ClickHouse")

		if err := etl.CreateTable(cfg.ClickHouseURL, ddl); err != nil {
			log.Error("failed to create table", zap.Error(err))
			return
		}

		log.Info("inserting data into ClickHouse")

		if err := etl.InsertRows(cfg.ClickHouseURL, cfg.Table, etl.GetColumnNames(td.Columns), td.Rows, cfg.BatchSize); err != nil {
			log.Error("failed to insert data", zap.Error(err))
			return
		}

		log.Success("initial data ingestion complete",
			zap.String("table", cfg.Table),
			zap.Int("rows", len(td.Rows)))

		if cfg.Polling.Enabled {
			ui.PrintSubtitle("Starting change data polling")

			lastSeen, err := determineLastSeen(td, cfg.Polling.Deltacol)

			if err != nil {
				log.Error("failed to determine last seen value", zap.Error(err))
				return
			}

			if err := startPolling(ctx, cfg, lastSeen); err != nil {
				log.Error("failed to start polling", zap.Error(err))
				return
			}

		}

		log.Success("data ingestion complete",
			zap.String("table", cfg.Table),
			zap.Int("rows", len(td.Rows)),
		)
	},
}

func loadConfig() *config.Config {
	log := log.StyledLog

	cfg, err := config.LoadConfig(ingestConfigPath)

	if err != nil {
		log.Warn("Could not load config from file, falling back to flags", zap.Error(err))
		cfg = &config.Config{
			PostgreSQLURL: ingestPgURL,
			ClickHouseURL: ingestChURL,
			Table:         ingestTable,
			Limit:         ingestLimit,
			BatchSize:     ingestBatch,
			Polling: config.PollingConfig{
				Enabled:  ingestPoll,
				Deltacol: ingestPollDelta,
				Interval: ingestPollInt,
			},
		}
	} else {
		if ingestPgURL != "" {
			cfg.PostgreSQLURL = ingestPgURL
		}
		if ingestChURL != "" {
			cfg.ClickHouseURL = ingestChURL
		}
		if ingestTable != "" {
			cfg.Table = ingestTable
		}
		if ingestLimit != 0 {
			cfg.Limit = ingestLimit
		}
		if ingestBatch != 0 {
			cfg.BatchSize = ingestBatch
		}

		if ingestPoll {
			cfg.Polling.Enabled = true
		}
		if ingestPollDelta != "" {
			cfg.Polling.Deltacol = ingestPollDelta
		}
		if ingestPollInt != 0 {
			cfg.Polling.Interval = ingestPollInt
		}
	}

	return cfg

}

func validateConfig(cfg *config.Config) bool {
	log := log.StyledLog

	if cfg.PostgreSQLURL == "" || cfg.ClickHouseURL == "" || cfg.Table == "" {
		log.Error("Missing required config values. Provide them in YAML or as flags.",
			zap.String("pg_url", cfg.PostgreSQLURL),
			zap.String("ch_url", cfg.ClickHouseURL),
			zap.String("table", cfg.Table),
		)
		return false
	}

	if cfg.Polling.Enabled {
		if cfg.Polling.Deltacol == "" {
			log.Error("Missing delta column for polling. Provide it in YAML or with --poll-delta flag.")
			return false
		}
		if cfg.Polling.Interval <= 0 {
			log.Error("Invalid polling interval. Must be greater than 0.")
			return false
		}
	}

	return true
}

func UI_itoa(n int) string {
	if n == 0 {
		return "all"
	}
	return fmt.Sprintf("%d", n)

}

func init() {
	ingestCmd.Flags().StringVar(&ingestConfigPath, "config", "", "Path to YAML config file (default: .pgtoch.yaml)")
	ingestCmd.Flags().StringVar(&ingestPgURL, "pg-url", "", "PostgreSQL connection URL")
	ingestCmd.Flags().StringVar(&ingestChURL, "ch-url", "", "ClickHouse connection URL")
	ingestCmd.Flags().StringVar(&ingestTable, "table", "", "Table name to ingest")
	ingestCmd.Flags().IntVar(&ingestLimit, "limit", 1000, "Limit rows to fetch from PG")
	ingestCmd.Flags().IntVar(&ingestBatch, "batch-size", 500, "Rows per ClickHouse insert")
	ingestCmd.Flags().BoolVar(&ingestPoll, "poll", false, "Continue polling for changes after initial ingest")
	ingestCmd.Flags().StringVar(&ingestPollDelta, "poll-delta", "", "Column name to track changes (usually a timestamp)")
	ingestCmd.Flags().IntVar(&ingestPollInt, "poll-interval", 0, "Polling interval in seconds")
	rootCmd.AddCommand(ingestCmd)
}
