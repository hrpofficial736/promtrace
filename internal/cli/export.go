package cli

import (
	"encoding/json"
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/spf13/cobra"
	"os"
)

var exportCommand *cobra.Command = &cobra.Command{
	Use:   "export",
	Short: "used to dump traces in std out",
	Long:  "used to dump traces in std out",
	RunE:  exportRun,
}

func exportRun(cmd *cobra.Command, args []string) error {

	cfg, err := config.Load()

	if err != nil {
		logger.Log.Error("error while loading config", "error", err)
		return err
	}

	st, err := store.NewSQLiteStore(cfg.DBPath)

	if err != nil {
		logger.Log.Error("error", "error in getting trace", err)
		return err
	}

	traces, err := st.ListAllTraces()

	if err != nil {
		logger.Log.Error("error", "error in getting trace", err)
		return err
	}

	encoder := json.NewEncoder(os.Stdout)

	for _, t := range traces {
		encoder.Encode(t)
	}

	fmt.Println(tui.RenderStatus(true, "Successfully exported traces"))

	return nil
}

var exportFormatFlag string

func init() {
	exportCommand.Flags().StringVar(&exportFormatFlag, "format", "jsonl", "output format(jsonl)")
	rootCmd.AddCommand(exportCommand)
}
