package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version string populated during build.
	Version = "v0.3.0"

	// RootCmd represents the base command when called without subcommands.
	RootCmd = &cobra.Command{
		Use:   "optivor",
		Short: "Optivor is an open-source image infrastructure framework.",
		Long: `Optivor provides image transformation runtime, driver interfaces, and deployment tooling
that let you run a production-grade image pipeline on top of object storage you already own.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			versionFlag, _ := cmd.Flags().GetBool("version")
			if versionFlag {
				fmt.Printf("optivor version %s\n", Version)
				return nil
			}
			return cmd.Help()
		},
	}
)

func init() {
	RootCmd.Flags().BoolP("version", "v", false, "print optivor version")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return RootCmd.Execute()
}
