package cmd

import (
	"context"
	ui "pgtoch/internal/UI"
	"pgtoch/internal/db"
	"pgtoch/internal/log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	chURL, pgURL string
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Test Postgres and CLickhouse connection",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintTitle("Testing db connections")
		ui.PrintSubtitle("checking connectivity to postgres and clickhouse")

		log := log.StyledLog

		log.Info("Testing postgres connection")
		conn, err := db.ConnectPostgres(pgURL)
		if err != nil {
			log.Error("Postgres connection failed", zap.Error(err))
			return
		}

		defer conn.Close(context.Background())
		log.Success("Postgres connection successful")

		log.Info("Testing clickhouse connection")
		chConn, err := db.ConnectClickhouse(chURL)
		if err != nil {
			log.Error("Clickhouse connection failed", zap.Error(err))
			return
		}

		defer chConn.Close()
		log.Success("Clickhouse connection successful")

		ui.PrintBox("Connection test successful", "All databases are reachable")
	},
}

func init() {
	connectCmd.Flags().StringVar(&pgURL, "pg-url", "", "PostgreSQL connection string")
	connectCmd.Flags().StringVar(&chURL, "ch-url", "", "ClickHouse connection URL")
	connectCmd.MarkFlagRequired("pg-url")
	connectCmd.MarkFlagRequired("ch-url")
	rootCmd.AddCommand(connectCmd)
}
