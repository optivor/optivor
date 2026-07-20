package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type DriverMetadata struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	OptivorAPI string `json:"optivor_api"`
	Path       string `json:"path,omitempty"`
}

type DriverRegistry struct {
	Drivers map[string]DriverMetadata `json:"drivers"`
}

var registryFilePath string

func getRegistryPath() (string, error) {
	if registryFilePath != "" {
		return registryFilePath, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to locate user home directory: %w", err)
	}
	dir := filepath.Join(home, ".config", "optivor")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}
	return filepath.Join(dir, "drivers.json"), nil
}

func loadRegistry() (*DriverRegistry, error) {
	path, err := getRegistryPath()
	if err != nil {
		return nil, err
	}

	reg := &DriverRegistry{Drivers: make(map[string]DriverMetadata)}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return reg, nil
		}
		return nil, fmt.Errorf("failed to read registry file: %w", err)
	}

	if err := json.Unmarshal(data, reg); err != nil {
		return nil, fmt.Errorf("failed to parse registry file: %w", err)
	}
	if reg.Drivers == nil {
		reg.Drivers = make(map[string]DriverMetadata)
	}
	return reg, nil
}

func saveRegistry(reg *DriverRegistry) error {
	path, err := getRegistryPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry data: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

var driverCmd = &cobra.Command{
	Use:   "driver",
	Short: "Manage external storage provider drivers",
	Long:  `The driver subcommand group allows installing, listing, removing, and inspecting external storage provider driver binaries.`,
}

var driverInstallCmd = &cobra.Command{
	Use:   "install <path-or-url>",
	Short: "Install and validate a storage driver binary",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		resolvedPath, err := resolveDriverSource(target)
		if err != nil {
			return fmt.Errorf("failed to resolve driver source %s: %w", target, err)
		}

		absPath, err := filepath.Abs(resolvedPath)
		if err != nil {
			return fmt.Errorf("invalid driver path: %w", err)
		}

		fi, err := os.Stat(absPath)
		if err != nil {
			return fmt.Errorf("driver binary not found at %s: %w", absPath, err)
		}
		if fi.IsDir() {
			return fmt.Errorf("driver path %s is a directory, expected executable file", absPath)
		}

		out, err := exec.Command(absPath, "--optivor-handshake").Output()
		if err != nil {
			return fmt.Errorf("driver handshake failed for %s: %w", absPath, err)
		}

		var meta DriverMetadata
		if err := json.Unmarshal(out, &meta); err != nil {
			return fmt.Errorf("invalid handshake JSON output from %s: %w", absPath, err)
		}
		if meta.Name == "" {
			return fmt.Errorf("driver handshake response missing 'name' field")
		}

		meta.Path = absPath

		reg, err := loadRegistry()
		if err != nil {
			return err
		}

		reg.Drivers[meta.Name] = meta
		if err := saveRegistry(reg); err != nil {
			return err
		}

		fmt.Printf("Driver '%s' (v%s) successfully installed.\n", meta.Name, meta.Version)
		return nil
	},
}

func resolveDriverSource(target string) (string, error) {
	if !strings.HasPrefix(target, "github:") && !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		return target, nil
	}

	// For remote sources, return placeholder or local path if file exists
	if _, err := os.Stat(target); err == nil {
		return target, nil
	}

	return target, nil
}

var driverListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered storage drivers",
	RunE: func(cmd *cobra.Command, args []string) error {
		reg, err := loadRegistry()
		if err != nil {
			return err
		}

		if len(reg.Drivers) == 0 {
			fmt.Println("No drivers installed.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tVERSION\tAPI\tPATH")
		for _, d := range reg.Drivers {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", d.Name, d.Version, d.OptivorAPI, d.Path)
		}
		return w.Flush()
	},
}

var driverRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a registered storage driver",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		reg, err := loadRegistry()
		if err != nil {
			return err
		}

		if _, exists := reg.Drivers[name]; !exists {
			return fmt.Errorf("driver '%s' is not installed", name)
		}

		delete(reg.Drivers, name)
		if err := saveRegistry(reg); err != nil {
			return err
		}

		fmt.Printf("Driver '%s' removed.\n", name)
		return nil
	},
}

var driverInfoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show details for a registered storage driver",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		reg, err := loadRegistry()
		if err != nil {
			return err
		}

		meta, exists := reg.Drivers[name]
		if !exists {
			return fmt.Errorf("driver '%s' is not installed", name)
		}

		fmt.Printf("Name:        %s\n", meta.Name)
		fmt.Printf("Version:     %s\n", meta.Version)
		fmt.Printf("Optivor API: %s\n", meta.OptivorAPI)
		fmt.Printf("Binary Path: %s\n", meta.Path)
		return nil
	},
}

func init() {
	driverCmd.AddCommand(driverInstallCmd)
	driverCmd.AddCommand(driverListCmd)
	driverCmd.AddCommand(driverRemoveCmd)
	driverCmd.AddCommand(driverInfoCmd)
	RootCmd.AddCommand(driverCmd)
}
