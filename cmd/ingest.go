package cmd

import ui "pgtoch/internal/UI"

var (
	ingestPgURL, ingestChURL, ingestTable, ingestConfigPath, ingestPollDelta string
	ingestLimit, ingestBatch, ingestPollInt                                  int
	ingestPoll                                                               bool
)

var ingestCmd = &cobra.Command{
	Use: "data ingest",
	Short: "transeferring postgres to clickhouse",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintTitle("Data Ingestion")
		ui.PrintSubtitle("transferring postgres to clickhouse")

		
}