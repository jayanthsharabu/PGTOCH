package cmd

import (
	"os"
	ui "pgtoch/internal/UI"
	"pgtoch/internal/log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var sampleConfigCmd = &cobra.Command{
	Use:   "sample-config",
	Short: "generate sample config file",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintTitle("Generating Sample Config")
		ui.PrintSubtitle("creating a sample config file")

		log := log.StyledLog
		log.Info("Generating Sample Config")

		const sampleConfig = `# pgtoch.yaml - Sample Config
# PostgreSQL connection URL
pg_url: "postgres://postgres:password@localhost:5432/mydb?sslmode=disable"

# ClickHouse HTTP interface URL
ch_url: "http://localhost:9000"

# Table to ingest from Postgres
table: UserAnswer

# Max rows to fetch
limit: 1000

# Batch size per insert
batch_size: 200

# Polling configuration
polling:
  # Enable polling for changes after initial ingest
  enabled: false
  # Column name to track changes (usually a timestamp)
  delta_column: "updated_at"
  # Polling interval in seconds
  interval_seconds: 30
`
		log.Info("Sample config generated successfully")
		err := os.WriteFile(".pgtoch.yaml", []byte(sampleConfig), 0644)
		if err != nil {
			log.Error("Failed to write sample config", zap.Error(err))
			return
		}

		log.Success("Sample config generated successfully")
		ui.PrintBox("Next Steps", "1. Edit the .pgtoch.yaml file with ur credentials\n"+"2. Configure ur table and polling setting"+"3. Run pgtoch ingest to start data transfer")
	},
}

func init() {
	rootCmd.AddCommand(sampleConfigCmd)
}
