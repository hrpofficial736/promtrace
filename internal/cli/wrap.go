package cli

import (
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/certmanager"
	"github.com/hrpofficial736/promtrace/internal/config"
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
	Short:                 "Wrap a subprocess and intercept its LLM API calls!",
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
		errString := tui.RenderStatus(false, "could not load configuration. please run setup and try again.")
		fmt.Println(errString)
		fmt.Println(tui.RenderStatus(false, "error while launching subprocess, please try again!"))
		return nil
	}
	// load CA
	cm, err := certmanager.NewCertManager(cfg)
	if err != nil {
		fmt.Println(tui.RenderStatus(false, "could not initialize certificate manager. please run setup and try again."))
		return nil
	}
	if err := cm.LoadCA(); err != nil {
		fmt.Println(tui.RenderStatus(false, "could not load root CA. run `promtrace setup` and try again."))
		return nil
	}

	// creating the store

	st, err := store.NewSQLiteStore(cfg.DBPath)

	if err != nil {
		fmt.Println(tui.RenderStatus(false, "could not open local trace store. please run setup and try again."))
		return nil
	}

	defer func() {
		err = st.Close()
		if err != nil {
			fmt.Println(tui.RenderStatus(false, "error closing the setup"))
		}
	}()

	sessionID := util.GenerateID()

	// start proxy server
	port := fmt.Sprintf(":%d", cfg.Proxy.Port)
	ps := proxy.NewServer(cm, st, port, sessionID)
	go func() {
		err = ps.StartServer()
		if err != nil {
			fmt.Println(tui.RenderStatus(false, "error starting the server"))
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// launch subprocess
	child, err := subprocess.LaunchChildProcessWithEnvVars(args, "127.0.0.1"+port, cm.GetCertPath())
	if err != nil {
		fmt.Println(tui.RenderStatus(false, "could not launch the wrapped command. check the command and try again."))
		return nil
	}

	err = child.Wait()
	if err != nil {
		fmt.Println(tui.RenderStatus(false, "error waiting for subprocess"))
	}
	return ps.Shutdown()
}

func init() {
	rootCmd.AddCommand(wrapCommand)
}
