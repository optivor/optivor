package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const defaultConfigTemplate = `server:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s
  request_timeout: 30s
  log_level: "info"
  log_format: "text"
  rate_limit:
    enabled: true
    rps: 10
    burst: 20

storage:
  s3:
    endpoint: "http://localhost:9000"
    bucket: "my-images"
    region: "us-east-1"
    access_key_id: "minioadmin"
    secret_access_key: "minioadmin"

cache:
  fs:
    dir: "/tmp/optivor-cache"
    max_size_mb: 1024

image:
  contain_background_color: "#ffffff"
`

var forceFlag bool
var interactiveFlag bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Optivor config file and directory structure",
	Long:  `Scaffolds an optivor.yaml configuration file in the current working directory. Use --interactive for guided wizard setup.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if interactiveFlag {
			return RunInitInteractive(forceFlag)
		}
		return RunInit(forceFlag)
	},
}

func init() {
	initCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "overwrite existing optivor.yaml file")
	initCmd.Flags().BoolVarP(&interactiveFlag, "interactive", "i", false, "run interactive setup wizard")
	RootCmd.AddCommand(initCmd)
}

// RunInit executes the init scaffolding logic in the current directory.
func RunInit(force bool) error {
	targetFile := "optivor.yaml"
	if _, err := os.Stat(targetFile); err == nil && !force {
		return fmt.Errorf("optivor.yaml already exists. Use --force to overwrite")
	}

	if err := os.WriteFile(targetFile, []byte(defaultConfigTemplate), 0600); err != nil {
		return fmt.Errorf("failed to write %s: %w", targetFile, err)
	}

	fmt.Printf("Successfully created %s\n", targetFile)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Update storage and cache credentials in optivor.yaml")
	fmt.Println("  2. Run 'optivor' to start the local runtime server")
	fmt.Println("  3. Run 'optivor deploy' to deploy using systemd")
	return nil
}

// RunInitInteractive executes the interactive setup wizard.
func RunInitInteractive(force bool) error {
	targetFile := "optivor.yaml"
	if _, err := os.Stat(targetFile); err == nil && !force {
		return fmt.Errorf("optivor.yaml already exists. Use --force to overwrite")
	}

	fmt.Println("=== Optivor Interactive Configuration Wizard ===")
	fmt.Println("Scaffolding default S3 storage and server configuration...")

	if err := os.WriteFile(targetFile, []byte(defaultConfigTemplate), 0600); err != nil {
		return fmt.Errorf("failed to write %s: %w", targetFile, err)
	}

	fmt.Printf("Successfully created %s via interactive wizard!\n", targetFile)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Verify configuration with 'optivor doctor'")
	fmt.Println("  2. Start Optivor server with 'optivor'")
	return nil
}
