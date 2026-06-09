package cli

import (
	"log/slog"
	"time"

	"github.com/hrpofficial736/promtrace/internal/certmanager"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/proxy"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/subprocess"
	"github.com/hrpofficial736/promtrace/internal/util"
	"github.com/spf13/cobra"
)

var wrapCommand *cobra.Command = &cobra.Command{
	Use:   "wrap",
	Short: "wrap a subprocess and intercept its LLM API calls!",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return wrapRun(cmd, args)
	},
}

func wrapRun(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	logger.Init(slog.LevelDebug)
	if err != nil {
		logger.Log.Error("error while loading config", "error", err)
		return err
	}
	logger.Log.Info("entered the wrap run function")
	// load CA
	cm, err := certmanager.NewCertManager(cfg)
	if err != nil {
		logger.Log.Error("error while making new cert manager", "error", err)
		return err
	}
	if err := cm.LoadCA(); err != nil {
		logger.Log.Error("error while loading CA", "error", err)
		return err
	}

	logger.Log.Info("yha tk aaya")

	// creating the store

	st, err := store.NewSQLiteStore(cm.GetDBPath())

	if err != nil {
		return err
	}
	logger.Log.Info("store bnne ke baad bhi hu")

	defer st.Close()

	sessionID := util.GenerateID()

	logger.Log.Info("id bhi generate ho gyi")

	// start proxy server
	ps := proxy.NewServer(cm, st, ":9117", sessionID)
	go ps.StartServer()

	time.Sleep(100 * time.Millisecond)

	// launch subprocess
	child, err := subprocess.LaunchChildProcessWithEnvVars(args, "127.0.0.1:9117", cm.GetCertPath())
	if err != nil {
		logger.Log.Error("error while launching child subprocess", "error", err)
		return err
	}

	err = child.Wait()

	ps.Shutdown()
	return nil
}

func init() {
	rootCmd.AddCommand(wrapCommand)
}
