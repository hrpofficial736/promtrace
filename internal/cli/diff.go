package cli

import (
	"fmt"
	"log/slog"

	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/diff"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/spf13/cobra"
)

var diffCommand *cobra.Command = &cobra.Command{
	Use:   "diff",
	Short: "used to comapare two traces",
	Long:  "used to compare two traces",
	Args:  cobra.ExactArgs(2),
	RunE:  diffRun,
}

func diffRun(cmd *cobra.Command, args []string) error {
	logger.Init(slog.LevelDebug)

	cfg, err := config.Load()

	if err != nil {
		logger.Log.Error("error while loading config", "error", err)
		return err
	}

	st, err := store.NewSQLiteStore(cfg.DBPath)

	if err != nil {
		return err
	}

	t1, _ := st.GetTrace(args[0])
	t2, _ := st.GetTrace(args[1])

	result, err := diff.DiffTraces(t1, t2)

	if err != nil {
		return err
	}

	fmt.Println(tui.RenderDiff(result))

	return nil
}

func init() {
	rootCmd.AddCommand(diffCommand)
}
