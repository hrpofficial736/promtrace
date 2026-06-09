package cli

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/spf13/cobra"
)

var exportCommand *cobra.Command = &cobra.Command{
	Use:   "export",
	Short: "used to dump traces in std out",
	Long:  "used to dump traces in std out",
	RunE:  exportRun,
}

func exportRun(cmd *cobra.Command, args []string) error {
	logger.Init(slog.LevelDebug)

	cfg, err := config.Load()

	if err != nil {
		logger.Log.Error("error while loading config", "error", err)
		return err
	}

	st, _ := store.NewSQLiteStore(cfg.DBPath)

	traces, err := st.ListAllTraces()

	if err != nil {
		logger.Log.Error("error", "error in getting trace", err)
		return nil
	}

	encoder := json.NewEncoder(os.Stdout)

	for _, t := range traces {
		encoder.Encode(t)
	}

	return nil
}

var exportFormatFlag string

func init() {
	exportCommand.Flags().StringVar(&exportFormatFlag, "format", "jsonl", "output format(jsonl)")
	rootCmd.AddCommand(exportCommand)
}
