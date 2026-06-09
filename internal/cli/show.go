package cli

import (
	"log/slog"

	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/spf13/cobra"
)

var showCommand *cobra.Command = &cobra.Command{
	Use:   "show",
	Short: "used to see a specific trace",
	Long:  "used to see a specific trace",
	RunE:  showRun,
}

func showRun(cmd *cobra.Command, args []string) error {
	logger.Init(slog.LevelDebug)
	cfg, err := config.Load()
	if err != nil {
		logger.Log.Error("error while loading config", "error", err)
		return err
	}

	st, _ := store.NewSQLiteStore(cfg.DBPath)

	trace, err := st.GetTrace(args[0])
	if err != nil {
		logger.Log.Error("error", "error in getting trace", err)
		return nil
	}

	logger.Log.Info("data", "trace is", trace)
	return nil
}

func init() {
	rootCmd.AddCommand(showCommand)
}
