package cli

import (
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/hrpofficial736/promtrace/internal/util"
	"github.com/spf13/cobra"
	"os"
)

var rootLongText string = tui.BoxWrapper(
	tui.RenderText(
		tui.Heading,
		util.GetPromtraceHeadingAsciiText()+
			"\n"+
			"promtrace intercepts LLM API calls from any process, logs prompts, responses, tokens, and cost — with zero code changes.",
	) + fmt.Sprintf("\n\nRun %s to get started.", tui.RenderText(tui.Hint, "'promtrace setup'")),
)

var rootCmd = &cobra.Command{
	Use:  "promtrace",
	Long: rootLongText,
}

func Execute() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(fmt.Errorf("error loading config: %s", err))
		return
	}

	logger.Init(cfg.Logging.Level)

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
