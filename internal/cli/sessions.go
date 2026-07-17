package cli

import (
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/spf13/cobra"
)

var sessionsLongText string = tui.BoxWrapper(
	tui.RenderCommandDescription(
		"sessions",
		"This command lists all recorded sessions with aggregate stats: total calls, cost, tokens, and average latency.",
	),
)
var sessionsCommand *cobra.Command = &cobra.Command{
	Use:                   "sessions",
	Short:                 "List all sessions with aggregate stats",
	DisableFlagsInUseLine: true,
	Long:                  sessionsLongText,
	RunE:                  sessionsRun,
}

func sessionsRun(cmd *cobra.Command, args []string) error {

	cfg, err := config.Load()
	if err != nil {
		logger.Log.Error("error loading config", "error", err)
		return err
	}
	store, err := store.NewSQLiteStore(cfg.DBPath)

	if err != nil {
		logger.Log.Error("error while making store", "error", err)
		return err
	}

	sessions, err := store.ListSessions()
	if err != nil {
		return err
	}

	fmt.Print(tui.RenderSessions(sessions))

	return nil
}

func init() {
	rootCmd.AddCommand(sessionsCommand)
}
