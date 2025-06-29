package cmd

import (
	"fmt"
	"path/filepath"
	"pgtoch/config"
	ui "pgtoch/internal/UI"
	"pgtoch/internal/etl/export"
	"pgtoch/internal/log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	exportTable, exportFormat, exportOut, exportChurl, exportConfigPath string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export from clickhouse ==> csv",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintTitle("Exporting data")
		ui.PrintSubtitle("clickhouse table to csv")

		log := log.StyledLog
		log.Info("Starting export..")

		cfg, err := config.LoadConfig(exportConfigPath)
		if err != nil {
			log.Error("Failed to load config falling to def", zap.Error(err))
			cfg = &config.Config{
				ClickHouseURL: exportChurl,
				Table:         exportTable,
			}
		} else {
			if exportChurl != "" {
				cfg.ClickHouseURL = exportChurl
			}
			if exportTable != "" {
				cfg.Table = exportTable
			}
		}

		if cfg.ClickHouseURL == "" || cfg.Table == "" {
			log.Fatal("missing required config values", zap.String("clickhouse_url", cfg.ClickHouseURL), zap.String("table", cfg.Table))
			return
		}

		ui.PrintBox("Export Details", fmt.Sprintf("Table: %s\nFormat: %s\nOutput: %s", cfg.Table, exportFormat, exportOut))

		outPath := filepath.Join(exportOut, fmt.Sprintf("%s.%s", cfg.Table, exportFormat))

		log.Info("Starting extraction")

		if err := export.ExportTabletoCSV(cfg.ClickHouseURL, cfg.Table, outPath); err != nil {
			log.Error("Failed to export ", zap.Error(err), zap.String("table", cfg.Table), zap.String("format", exportFormat), zap.String("output", outPath))
			return
		}

		log.Success("Export completed successfully")
		ui.PrintBox("Export Complete", fmt.Sprintf("table %s format %s output file : %s", cfg.Table, exportFormat, outPath))

	},
}

func init() {
	exportCmd.Flags().StringVar(&exportConfigPath, "config", "", "Path to YAML config file")
	exportCmd.Flags().StringVar(&exportChurl, "ch-url", "", "ClickHouse connection URL")
	exportCmd.Flags().StringVar(&exportTable, "table", "", "Table name to export")
	exportCmd.Flags().StringVar(&exportFormat, "format", "csv", "Export format (csv)")
	exportCmd.Flags().StringVar(&exportOut, "out", ".", "Output directory for exported files")

	exportCmd.MarkFlagRequired("format")
	exportCmd.MarkFlagRequired("out")

	rootCmd.AddCommand(exportCmd)
}
