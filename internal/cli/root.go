package cli

import (
	"log/slog"
	"os"

	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "promtrace",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	logger.Init(slog.LevelDebug)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
