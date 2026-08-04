package cli

import (
	"encoding/json"
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/spf13/cobra"
	"os"
)

var exportLongText string = tui.BoxWrapper(
	tui.RenderCommandDescription(
		"export",
		"This command dumps all the traces to the standard output in JSONL format by default.",
	),
)
var exportCommand *cobra.Command = &cobra.Command{
	Use:                   "export",
	Example:               "  promtrace export --format jsonl > traces.jsonl",
	Short:                 "Export all traces as newline-delimited JSON",
	DisableFlagsInUseLine: true,
	Long:                  exportLongText,
	RunE:                  exportRun,
}

func exportRun(cmd *cobra.Command, args []string) error {

	cfg, err := config.Load()

	if err != nil {
		errString := tui.RenderStatus(false, "could not load configuration. please run setup and try again.")
		fmt.Println(errString)
		return nil
	}

	st, err := store.NewSQLiteStore(cfg.DBPath)

	if err != nil {
		errString := tui.RenderStatus(false, "error occured, please try again")
		fmt.Println(errString)
		return nil
	}

	traces, err := st.ListAllTraces()

	if err != nil {
		errString := tui.RenderStatus(false, "error occured, please try again")
		fmt.Println(errString)
		return nil
	}

	encoder := json.NewEncoder(os.Stdout)

	switch exportFormatFlag {
	case "jsonl":
		for _, t := range traces {
			if err := encoder.Encode(t); err != nil {
				err := encoder.Encode(t)
				if err != nil {
					errString := tui.RenderStatus(false, "error occured, please try again")
					fmt.Println(errString)
					return nil
				}
			}
		}

	case "json":
		encoder.SetIndent("", "    ")
		return encoder.Encode(traces)

	default:
		errString := tui.RenderStatus(false, fmt.Sprintf("unsupported format %q (supported: jsonl, json)", exportFormatFlag))
		fmt.Println(errString)
		return nil
	}

	return nil
}

var exportFormatFlag string

func init() {
	exportCommand.Flags().StringVar(&exportFormatFlag, "format", "jsonl", "output format (jsonl, json)")
	rootCmd.AddCommand(exportCommand)
}
