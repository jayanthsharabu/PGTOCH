package cmd

import (
	"fmt"
	"os"
	ui "pgtoch/internal/UI"
	"pgtoch/internal/log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	useInteractive bool
)

var rootCmd = &cobra.Command{
	Use:   "pgtoch",
	Short: "Etl from postgres ==> clickhouse",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if useInteractive && cmd.Name() == "pgtoch" {
			showInteractiveUI()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		showLogo()
	},
}

func showLogo() {
	ui.PrintLogo()
	ui.PrintTitle("pgtoch")
	ui.PrintSubtitle("Etl from postgres ==> clickhouse")
	fmt.Println()
}

func Execute() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "help" || os.Args[1] == "-h") {
		showLogo()
	}
	if err := rootCmd.Execute(); err != nil {
		log.StyledLog.Error("Error executing command", zap.Error(err))
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&useInteractive, "interactive", "i", false, "Use Interactive mode TUI Mode")
}

func showInteractiveUI() {
	p := tea.NewProgram(ui.NewAppModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.StyledLog.Error("Error running interactive UI", zap.Error(err))
		os.Exit(1)
	}
	os.Exit(0)
}
