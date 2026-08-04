package cli

import (
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/config"
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
		errString := tui.RenderStatus(false, "could not load configuration. please run setup and try again.")
		fmt.Println(errString)
		return nil
	}
	store, err := store.NewSQLiteStore(cfg.DBPath)

	if err != nil {
		fmt.Println(tui.RenderStatus(false, "could not open local trace store. please run setup and try again."))
		return nil
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
