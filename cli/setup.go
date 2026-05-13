package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var setupCommand *cobra.Command = &cobra.Command{
	Use:   "setup",
	Short: "setup tls",
	Long:  "setup tls",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("setup called")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCommand)
}
