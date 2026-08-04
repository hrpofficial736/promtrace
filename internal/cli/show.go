package cli

import (
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/spf13/cobra"
)

var showLongText string = tui.BoxWrapper(
	tui.RenderCommandDescription(
		"show",
		"This command displays the full details of an intercepted trace including model, prompts, response, cost, tokens, and latency.",
	),
)
var showCommand *cobra.Command = &cobra.Command{
	Use:                   "show <TRACE_ID>",
	Example:               "  promtrace show a3f7b2c1",
	Short:                 "Display details of a specific trace",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	Long:                  showLongText,
	RunE:                  showRun,
}

func showRun(cmd *cobra.Command, args []string) error {
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

	trace, err := st.GetTrace(args[0])
	if err != nil {
		fmt.Println(tui.RenderStatus(false, "trace not found. please check the trace ID and try again."))
		return nil
	}

	fmt.Print(tui.RenderTraceInfoContainer(trace))

	return nil
}

func init() {
	rootCmd.AddCommand(showCommand)
}
