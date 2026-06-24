package cli

import (
	"fmt"
	"log/slog"

	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/spf13/cobra"
)

var sessionsCommand *cobra.Command = &cobra.Command{
	Use:   "sessions",
	Short: "used to list all the sessions in this device",
	Long:  "used to list all the sessions in this device",
	RunE:  sessionsRun,
}

func sessionsRun(cmd *cobra.Command, args []string) error {
	logger.Init(slog.LevelDebug)

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

	fmt.Println(tui.RenderSessions(sessions))

	return nil
}

func init() {
	rootCmd.AddCommand(sessionsCommand)
}
