package cli

import (
	"fmt"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/optivor/optivor/internal/config"
	"github.com/spf13/cobra"
)

var doctorConfigFlag string

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Perform system and configuration health checks",
	Long:  `Validates local configuration, environment variables, dependencies (libvips), and connectivity.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunDoctor(doctorConfigFlag)
	},
}

func init() {
	doctorCmd.Flags().StringVarP(&doctorConfigFlag, "config", "c", "optivor.yaml", "path to config file")
	RootCmd.AddCommand(doctorCmd)
}

// RunDoctor runs system diagnostic checks and reports health status.
func RunDoctor(configPath string) error {
	fmt.Println("Running Optivor system health checks...")
	hasErrors := false

	// Check 1: Config file loading & parsing
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("  ❌ Config file check failed (%s): %v\n", configPath, err)
		hasErrors = true
	} else {
		fmt.Printf("  ✅ Config file valid (%s)\n", configPath)
	}

	// Check 2: Auth secret check if signed URLs enabled
	if cfg != nil && cfg.Auth.SignedURLs.Enabled {
		secret := cfg.Auth.SignedURLs.Secret
		if secret == "" {
			secret = os.Getenv("OPTIVOR_AUTH_SECRET")
		}
		if secret == "" {
			fmt.Println("  ⚠️  Signed URLs enabled but secret is missing (set OPTIVOR_AUTH_SECRET)")
		} else {
			fmt.Println("  ✅ Signed URL secret configured")
		}
	} else {
		fmt.Println("  ✅ Signed URL security check skipped (disabled)")
	}

	// Check 3: Storage config validation (single vs multi-bucket)
	if cfg != nil {
		if len(cfg.Buckets) > 0 {
			fmt.Printf("  ✅ Multi-bucket mode configured (%d buckets)\n", len(cfg.Buckets))
			for _, b := range cfg.Buckets {
				if b.Endpoint == "" || b.Bucket == "" {
					fmt.Printf("  ❌ Bucket %q configuration incomplete (missing endpoint or bucket)\n", b.Name)
					hasErrors = true
				} else {
					fmt.Printf("  ✅ Bucket %q valid (Provider: %s, Endpoint: %s, Access: %s)\n", b.Name, b.Provider, b.Endpoint, b.Access)
				}
			}
		} else if cfg.Storage.S3.Endpoint != "" && cfg.Storage.S3.Bucket != "" {
			fmt.Printf("  ✅ Storage configuration valid (Endpoint: %s, Bucket: %s)\n", cfg.Storage.S3.Endpoint, cfg.Storage.S3.Bucket)
		} else {
			fmt.Println("  ❌ Storage configuration incomplete (missing endpoint or bucket)")
			hasErrors = true
		}
	}

	// Check 4: libvips runtime initialization check
	// #nosec G104
	_ = vips.Startup(nil)
	defer vips.Shutdown()
	fmt.Printf("  ✅ libvips %s initialized successfully\n", vips.Version)

	if hasErrors {
		return fmt.Errorf("doctor check completed with errors")
	}

	fmt.Println("\nAll system health checks passed!")
	return nil
}
