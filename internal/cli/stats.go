package cli

import (
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/tui"
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

	cfg, err := config.Load()

	if err != nil {
		logger.Log.Error("error while loading config", "error", err)
		return err
	}

	st, err := store.NewSQLiteStore(cfg.DBPath)

	if err != nil {
		logger.Log.Error("error while parsing duration", "error", err)
		return err
	}

	duration, err := util.ParseDuration(statsFormatFlag)

	if err != nil {
		logger.Log.Error("error while parsing duration", "error", err)
		return err
	}

	stats, err := st.GetStats(duration)

	if err != nil {
		logger.Log.Error("error", "error in getting stats", err)
		return err
	}

	fmt.Println(tui.RenderStats(stats))

	return nil
}

var statsFormatFlag string

func init() {
	statsCommand.Flags().StringVar(&statsFormatFlag, "last", "7d", "time window (e.g., 7d, 24h, 30m)")
	rootCmd.AddCommand(statsCommand)
}
