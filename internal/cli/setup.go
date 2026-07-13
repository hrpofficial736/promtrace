package cli

import (
	"fmt"
	"github.com/hrpofficial736/promtrace/internal/certmanager"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/spf13/cobra"
)

var setupLongText string = tui.RenderCommandDescription("setup", "This command generates root CA and installs it in the OS trust store.")

var setupCommand *cobra.Command = &cobra.Command{
	Use:                   "setup",
	Short:                 "setup tls",
	Long:                  setupLongText,
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()

		if err != nil {
			logger.Log.Error("error while loading config", "error", err)
			fmt.Println(tui.RenderStatus(false, "Error generating certificates, please try again!"))
			return err
		}
		cm, err := certmanager.NewCertManager(cfg)
		if err != nil {
			logger.Log.Info("error while generating cert manager")
			fmt.Println(tui.RenderStatus(false, "Error generating certificates, please try again!"))
			return err
		}

		if err := cm.GenerateRootCACertificate(); err != nil {
			logger.Log.Info("error while generating ca cert")
			fmt.Println(tui.RenderStatus(false, "Error generating certificates, please try again!"))
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
