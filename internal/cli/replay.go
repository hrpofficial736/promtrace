package cli

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/provider"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/hrpofficial736/promtrace/internal/util"
	"github.com/hrpofficial736/promtrace/pkg/costable"
	"github.com/spf13/cobra"
)

var replayLongText string = tui.BoxWrapper(
	tui.RenderCommandDescription(
		"replay",
		"This command re-sends a captured LLM request, optionally with a different model from the same provider. Saves the result as a new trace and compares it against the original.",
	),
)
var replayCommand *cobra.Command = &cobra.Command{
	Use:                   "replay <TRACE_ID> --model <MODEL_NAME>",
	Example:               "  promtrace replay a3f7b2c1 --model gpt-4o-mini",
	Short:                 "Re-send a captured LLM call with an optional model swap",
	Long:                  replayLongText,
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	RunE:                  replayRun,
}

func replayRun(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()

	if err != nil {
		errString := tui.RenderStatus(false, "could not load configuration. please run setup and try again.")
		fmt.Println(errString)
		return nil
	}

	st, err := store.NewSQLiteStore(cfg.DBPath)

	if err != nil {
		errString := tui.RenderStatus(false, "could not open local trace store. please run setup and try again.")
		fmt.Println(errString)
		return nil
	}

	t, err := st.GetTrace(args[0])

	if err != nil {
		fmt.Println(tui.RenderStatus(false, "trace not found. please check the trace ID and try again."))
		return nil
	}

	if replayModelFlag != "" {
		err = provider.ValidateModel(replayModelFlag, t.Host)
		if err != nil {
			fmt.Println(tui.RenderStatus(false, "could not replay the request. the selected model is either unsupported or does not match the original provider family. please choose a compatible model and try again."))
			return nil
		}
	}

	start := time.Now()

	m := tui.NewSpinnerModel(t, replayModelFlag)

	result, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println(tui.RenderStatus(false, "replay failed to start. please try again."))
		return nil
	}

	finalModel := result.(tui.SpinnerModel)

	if finalModel.Err != nil {
		fmt.Println(tui.RenderStatus(false, "could not replay the request. please try again."))
		return nil
	}

	resBytes := finalModel.Body

	latency := time.Since(start).Milliseconds()

	newTrace := *t

	newTrace.ID = util.GenerateID()

	newTrace.SessionID = t.SessionID

	newTrace.Timestamp = time.Now()

	extractorRes, inTokens, outTokens := provider.GetExtractor(t.Host).ExtractResponse(resBytes)

	newTrace.Tokens = inTokens + outTokens

	newTrace.Response = extractorRes

	newTrace.StatusCode = finalModel.Res.StatusCode

	newTrace.LatencyMs = latency

	if replayModelFlag != "" {
		newTrace.Model = replayModelFlag
	}

	newTrace.Cost = costable.CalculateCost(newTrace.Model, inTokens, outTokens)

	err = st.SaveTrace(&newTrace)

	if err != nil {
		fmt.Println(tui.RenderStatus(false, "could not save replay result. please try again."))
		return nil
	}

	r1 := &tui.ReplayComparisonResponseStruct{
		Model:          t.Model,
		Latency:        int(t.LatencyMs),
		Cost:           t.Cost,
		ResponseLength: len(t.Response),
	}

	r2 := &tui.ReplayComparisonResponseStruct{
		Model:          newTrace.Model,
		Latency:        int(newTrace.LatencyMs),
		Cost:           newTrace.Cost,
		ResponseLength: len(newTrace.Response),
	}

	fmt.Println(tui.RenderReplay(r1, r2))

	return nil

}

var replayModelFlag string

func init() {
	replayCommand.Flags().StringVar(&replayModelFlag, "model", "", "model to use for replay (must be from same provider)")
	rootCmd.AddCommand(replayCommand)
}
