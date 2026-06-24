package cli

import (
	"fmt"
	"log/slog"

	"github.com/hrpofficial736/promtrace/internal/certmanager"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/tui"
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
			fmt.Println(tui.RenderStatus(false, "Error generatinc certificates, please try again!"))
			return err
		}
		fmt.Println("setup called")
		cm, err := certmanager.NewCertManager(cfg)
		if err != nil {
			logger.Log.Info("error while generating cert manager")
			fmt.Println(tui.RenderStatus(false, "Error generatinc certificates, please try again!"))
			return err
		}

		if err := cm.GenerateRootCACertificate(); err != nil {
			logger.Log.Info("error while generating ca cert")
			fmt.Println(tui.RenderStatus(false, "Error generatinc certificates, please try again!"))
			return err
		}

		cm.StoreRootCACertificateInSystemTrustStore()
		// setup complete, promtrace is ready

		fmt.Println(tui.RenderStatus(true, "Certificates setup is complete, promtrace is ready."))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCommand)
}
