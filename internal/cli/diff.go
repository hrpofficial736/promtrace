package cli

import (
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/diff"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/spf13/cobra"
)

var diffLongText string = tui.BoxWrapper(
	tui.RenderCommandDescription(
		"diff",
		"This command compares two traces by their system prompts, user prompts, cost, tokens and latency(ms).",
	),
)
var diffCommand *cobra.Command = &cobra.Command{
	Use:                   "diff <TRACE_1_ID> <TRACE_2_ID>",
	Example:               "  promtrace diff a3f7b2c1 b8c4d1e2",
	Short:                 "Compare two traces side by side",
	Long:                  diffLongText,
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(2),
	RunE:                  diffRun,
}

func diffRun(cmd *cobra.Command, args []string) error {

	cfg, err := config.Load()

	if err != nil {
		errString := tui.RenderStatus(false, "could not load configuration. please run setup and try again.")
		fmt.Println(errString)
		return nil
	}

	st, err := store.NewSQLiteStore(cfg.DBPath)

	if err != nil {
		return err
	}

	t1, err := st.GetTrace(args[0])
	if err != nil {
		errString := tui.RenderStatus(false, fmt.Sprintf("trace %q not found", args[0]))
		fmt.Println(errString)
		return nil
	}

	t2, err := st.GetTrace(args[1])
	if err != nil {
		errString := tui.RenderStatus(false, fmt.Sprintf("trace %q not found", args[1]))
		fmt.Println(errString)
		return nil
	}

	result, err := diff.DiffTraces(t1, t2)

	if err != nil {
		return err
	}

	fmt.Print(tui.RenderDiff(result))

	return nil
}

func init() {
	rootCmd.AddCommand(diffCommand)
}
