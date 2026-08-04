package cli

import (
	"fmt"

	"github.com/hrpofficial736/promtrace/internal/certmanager"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/tui"
	"github.com/spf13/cobra"
)

var setupLongText string = tui.BoxWrapper(
	tui.RenderCommandDescription(
		"setup",
		"This command generates root CA and installs it in the OS's trust store.",
	),
)

var setupCommand *cobra.Command = &cobra.Command{
	Use:                   "setup",
	Short:                 "Sets up promtrace for TLS interception.",
	Long:                  setupLongText,
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()

		if err != nil {
			errString := tui.RenderStatus(false, "could not load configuration. please try again.")
			fmt.Println(errString)
			return nil
		}
		cm, err := certmanager.NewCertManager(cfg)
		if err != nil {
			fmt.Println(tui.RenderStatus(false, "could not initialize certificate manager. please try again."))
			return nil
		}

		if err := cm.GenerateRootCACertificate(); err != nil {
			fmt.Println(tui.RenderStatus(false, "could not generate root CA certificate. please try again."))
			return nil
		}

		cm.StoreRootCACertificateInSystemTrustStore()
		// setup complete, promtrace is ready

		fmt.Println(tui.RenderStatus(true, "promtrace is ready, get started by wrapping a sub-process."))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCommand)
}
