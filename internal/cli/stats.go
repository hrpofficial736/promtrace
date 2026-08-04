package cli

import (
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/hrpofficial736/promtrace/internal/util"
	"github.com/spf13/cobra"
)

var statsLongText string = tui.BoxWrapper(
	tui.RenderCommandDescription(
		"stats",
		"This command displays the stats of the LLM usage since the specified duration.",
	),
)
var statsCommand *cobra.Command = &cobra.Command{
	Use:     "stats",
	Example: "  promtrace stats --last 7d",
	Short:   "Show cost and token trends over time",
	Long:    statsLongText,
	RunE:    statsRun,
}

func statsRun(cmd *cobra.Command, args []string) error {

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

	duration, err := util.ParseDuration(statsFormatFlag)

	if err != nil {
		fmt.Println(tui.RenderStatus(false, "please enter a valid duration (e.g., 7d, 2h, 40m, 50s)"))
		return nil
	}

	stats, err := st.GetStats(duration)

	if err != nil {
		fmt.Println(tui.RenderStatus(false, "could not fetch usage stats. please try again."))
		return nil
	}

	fmt.Println(tui.RenderStats(stats))

	return nil
}

var statsFormatFlag string

func init() {
	statsCommand.Flags().StringVar(&statsFormatFlag, "last", "7d", "time window (e.g., 7d, 24h, 30m)")
	rootCmd.AddCommand(statsCommand)
}
