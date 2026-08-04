package cli

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hrpofficial736/promtrace/internal/config"
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
		errString := tui.RenderStatus(false, "could not load configuration. please run setup and try again.")
		fmt.Println(errString)
		return nil
	}

	st, err := store.NewSQLiteStore(cfg.DBPath)

	if err != nil {
		fmt.Println(tui.RenderStatus(false, "could not open local trace store. please run setup and try again."))
		return nil
	}

	m := tui.NewWatchModel(st, cfg.Watch.Limit)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println(tui.RenderStatus(false, "watch mode failed to start. please try again."))
		return nil
	}

	return nil
}

func init() {
	rootCmd.AddCommand(watchCommand)
}
