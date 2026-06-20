package cli

import (
	"log/slog"

	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/spf13/cobra"
)

var watchCommand *cobra.Command = &cobra.Command{
	Use:   "watch",
	Short: "used to display live traces",
	Long:  "used to display live traces",
	RunE:  watchRun,
}

func watchRun(cmd *cobra.Command, args []string) error {
	logger.Init(slog.LevelDebug)

	_, err := config.Load()

	if err != nil {
		logger.Log.Error("error while loading config", "error", err)
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(watchCommand)
}
