package cli

import (
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/spf13/cobra"
)

var showCommand *cobra.Command = &cobra.Command{
	Use:   "show",
	Short: "used to see a specific trace",
	Long:  "used to see a specific trace",
	RunE:  showRun,
}

func showRun(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		logger.Log.Error("error while loading config", "error", err)
		return err
	}

	st, err := store.NewSQLiteStore(cfg.DBPath)

	if err != nil {
		logger.Log.Error("error while initializing store", "error", err)
		return err
	}

	trace, err := st.GetTrace(args[0])
	if err != nil {
		logger.Log.Error("error", "error in getting trace", err)
		return nil
	}

	fmt.Println(tui.RenderTraceInfoContainer(trace))

	return nil
}

func init() {
	rootCmd.AddCommand(showCommand)
}
