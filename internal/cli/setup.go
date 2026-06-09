package cli

import (
	"fmt"
	"log/slog"

	"github.com/hrpofficial736/promtrace/internal/certmanager"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/spf13/cobra"
)

var setupCommand *cobra.Command = &cobra.Command{
	Use:   "setup",
	Short: "setup tls",
	Long:  "setup tls",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Init(slog.LevelDebug)
		cfg, err := config.Load()

		if err != nil {
			logger.Log.Error("error while loading config", "error", err)
			return err
		}
		fmt.Println("setup called")
		cm, err := certmanager.NewCertManager(cfg)
		if err != nil {
			logger.Log.Info("error while generating cert manager")
			return err
		}

		if err := cm.GenerateRootCACertificate(); err != nil {
			logger.Log.Info("error while generating ca cert")
			return err
		}

		cm.StoreRootCACertificateInSystemTrustStore()
		// setup complete, promtrace is ready
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCommand)
}
