package cli

import (
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/provider"
	"github.com/hrpofficial736/promtrace/internal/replay"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/hrpofficial736/promtrace/internal/util"
	"github.com/spf13/cobra"
	"io"
	"time"
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
		logger.Log.Error("error while loading config", "error", err)
		return err
	}

	st, err := store.NewSQLiteStore(cfg.DBPath)

	if err != nil {
		logger.Log.Error("error while making store", "error", err)
		return err
	}

	t, err := st.GetTrace(args[0])

	if err != nil {
		logger.Log.Error("error while getting trace", "error", err)
		return err
	}

	if replayModelFlag != "" {
		err = provider.ValidateModel(replayModelFlag, t.Host)
		if err != nil {
			return err
		}

	}

	start := time.Now()

	res, err := replay.ReplayRequest(t, replayModelFlag)

	if err != nil {
		logger.Log.Error("error while re-sending the request", "error", err)
		return err
	}

	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)

	if err != nil {
		logger.Log.Error("error reading response body", "error", err)
		return err
	}

	newTrace := *t

	newTrace.ID = util.GenerateID()

	newTrace.SessionID = t.SessionID

	newTrace.Timestamp = time.Now()

	extractorRes, inTokens, outTokens := provider.GetExtractor(t.Host).ExtractResponse(resBytes)

	newTrace.Tokens = inTokens + outTokens

	newTrace.Response = extractorRes

	newTrace.StatusCode = res.StatusCode

	newTrace.LatencyMs = time.Since(start).Milliseconds()

	if replayModelFlag != "" {
		newTrace.Model = replayModelFlag
	}

	err = st.SaveTrace(&newTrace)

	if err != nil {
		logger.Log.Error("error saving replay trace", "error", err)
		return err
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
