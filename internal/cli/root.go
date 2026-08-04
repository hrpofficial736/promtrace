package cli

import (
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/config"
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
			"A transparent LLM call interceptor. See exactly what prompts your app sends, what they cost, and when they silently change — with zero code modifications.",
	) + fmt.Sprintf("\n\nRun %s to get started.", tui.RenderText(tui.Hint, "'promtrace setup'")),
)

var rootCmd = &cobra.Command{
	Use:  "promtrace",
	Long: rootLongText,
}

func Execute() {
	_, err := config.Load()
	if err != nil {
		errString := tui.RenderStatus(false, "could not load configuration. please run setup and try again.")
		fmt.Println(errString)
		return
	}

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
