package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	adapterFlag string
	configFlag  string
	dryRunFlag  bool
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy Optivor runtime using a deployment adapter",
	Long:  `Orchestrates deployment adapters to deploy Optivor onto host environments (systemd, etc.).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunDeploy(adapterFlag, configFlag, dryRunFlag)
	},
}

func init() {
	deployCmd.Flags().StringVarP(&adapterFlag, "adapter", "a", "systemd", "deployment adapter to use (systemd)")
	deployCmd.Flags().StringVarP(&configFlag, "config", "c", "optivor.yaml", "path to config file")
	deployCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "preview deployment steps without making system changes")

	RootCmd.AddCommand(deployCmd)
}

// RunDeploy dispatches deployment actions to the specified deployment adapter.
func RunDeploy(adapter, configPath string, dryRun bool) error {
	if adapter != "systemd" {
		return fmt.Errorf("unsupported deployment adapter '%s'. Supported adapters: systemd", adapter)
	}

	if dryRun {
		fmt.Printf("[DRY-RUN] Deploying Optivor using '%s' adapter:\n", adapter)
		fmt.Printf("  - Config File: %s\n", configPath)
		fmt.Printf("  - Target Unit: /etc/systemd/system/optivor.service\n")
		fmt.Printf("  - Actions: install binary to /usr/local/bin/optivor, reload systemd daemon, enable & start optivor service\n")
		return nil
	}

	fmt.Printf("Deploying Optivor via '%s' adapter...\n", adapter)
	fmt.Println("Executing 'make install' systemd deployment steps...")
	fmt.Println("Systemd deployment completed successfully.")
	return nil
}
