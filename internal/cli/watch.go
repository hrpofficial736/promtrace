package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/spf13/cobra"
)

var watchLongText string = tui.BoxWrapper(
	tui.RenderCommandDescription(
		"watch",
		"This command displays all the intercepted traces within current session in real time as the requests arrive.",
	),
)
var watchCommand *cobra.Command = &cobra.Command{
	Use:                   "watch",
	Short:                 "Live trace feed — updates in real time",
	DisableFlagsInUseLine: true,
	Long:                  watchLongText,
	RunE:                  watchRun,
}

func watchRun(cmd *cobra.Command, args []string) error {

	cfg, err := config.Load()

	if err != nil {
		logger.Log.Error("error while loading config", "error", err)
		return err
	}

	st, err := store.NewSQLiteStore(cfg.DBPath)

	if err != nil {
		logger.Log.Error("error", "error", err)
		return nil
	}

	m := tui.NewWatchModel(st, cfg.Watch.Limit)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		logger.Log.Error("error", "error", err)
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(watchCommand)
}
