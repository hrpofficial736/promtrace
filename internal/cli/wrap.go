package cli

import (
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/certmanager"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/proxy"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/subprocess"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/hrpofficial736/promtrace/internal/util"
	"github.com/spf13/cobra"
	"time"
)

var wrapLongText string = tui.BoxWrapper(
	tui.RenderCommandDescription(
		"wrap",
		"This command wraps a sub-process and starts intercepting all of the HTTP requests made by that sub-process.",
	),
)

var wrapCommand *cobra.Command = &cobra.Command{
	Use:                   "wrap <PROCESS_COMMAND>",
	Example:               "  promtrace wrap python script.py",
	Short:                 "wrap a subprocess and intercept its LLM API calls!",
	Long:                  wrapLongText,
	Args:                  cobra.MinimumNArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return wrapRun(cmd, args)
	},
}

func wrapRun(_ *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		logger.Log.Error("error while loading config", "error", err)
		fmt.Println(tui.RenderStatus(false, "error while launching subprocess, please try again!"))
		return err
	}
	// load CA
	cm, err := certmanager.NewCertManager(cfg)
	if err != nil {
		logger.Log.Error("error while making new cert manager", "error", err)
		fmt.Println(tui.RenderStatus(false, "error while launching subprocess, please try again!"))
		return err
	}
	if err := cm.LoadCA(); err != nil {
		logger.Log.Error("error while loading CA", "error", err)
		fmt.Println(tui.RenderStatus(false, "error while launching subprocess, please try again!"))
		return err
	}

	// creating the store

	st, err := store.NewSQLiteStore(cfg.DBPath)

	if err != nil {
		fmt.Println(tui.RenderStatus(false, "error while launching subprocess, please try again!"))
		return err
	}

	defer st.Close()

	sessionID := util.GenerateID()

	// start proxy server
	port := fmt.Sprintf(":%d", cfg.Proxy.Port)
	ps := proxy.NewServer(cm, st, port, sessionID)
	go ps.StartServer()

	time.Sleep(100 * time.Millisecond)

	// launch subprocess
	child, err := subprocess.LaunchChildProcessWithEnvVars(args, "127.0.0.1"+port, cm.GetCertPath())
	if err != nil {
		logger.Log.Error("error while launching child subprocess", "error", err)
		fmt.Println(tui.RenderStatus(false, "error while launching subprocess, please try again!"))
		return err
	}

	err = child.Wait()
	return ps.Shutdown()
}

func init() {
	rootCmd.AddCommand(wrapCommand)
}
