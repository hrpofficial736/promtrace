package cli

import (
	"log/slog"

	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/util"
	"github.com/spf13/cobra"
)

var statsCommand *cobra.Command = &cobra.Command{
	Use:   "stats",
	Short: "shows stats of the llm usage",
	Long:  "shows stats of the llm usage",
	RunE:  statsRun,
}

func statsRun(cmd *cobra.Command, args []string) error {
	logger.Init(slog.LevelDebug)

	cfg, err := config.Load()

	if err != nil {
		logger.Log.Error("error while loading config", "error", err)
		return err
	}

	st, _ := store.NewSQLiteStore(cfg.DBPath)

	duration, err := util.ParseDuration(statsFormatFlag)

	if err != nil {
		logger.Log.Error("error while parsing duration", "error", err)
		return err
	}

	stats, err := st.GetStats(duration)

	if err != nil {
		logger.Log.Error("error", "error in getting stats", err)
		return nil
	}

	logger.Log.Info("stats of last "+statsFormatFlag, "stats", stats)
	return nil
}

var statsFormatFlag string

func init() {
	statsCommand.Flags().StringVar(&statsFormatFlag, "last", "7d", "stats")
	rootCmd.AddCommand(statsCommand)
}
