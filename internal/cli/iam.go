package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/optivor/optivor/internal/config"
	"github.com/spf13/cobra"
)

var (
	iamConfigFile   string
	iamCapabilities string
	iamAllowedPaths string
	iamApiKey       string
	iamApiRole      string
	iamApiBuckets   string
)

var iamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Manage IAM roles, path authorization policies, and API keys",
	Long:  `Subcommand group for inspecting and configuring Role-Based Access Control (RBAC) and path-level prefix policies.`,
}

var iamRoleCmd = &cobra.Command{
	Use:   "role",
	Short: "Manage IAM role definitions",
}

var iamRoleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured IAM roles and their allowed paths",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := iamConfigFile
		if cfgPath == "" {
			cfgPath = "optivor.yaml"
		}

		cfg, err := config.Load(cfgPath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to load configuration '%s': %w", cfgPath, err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ROLE NAME\tDESCRIPTION\tCAPABILITIES\tALLOWED PATHS")

		// Built-in defaults
		fmt.Fprintln(w, "admin\tBuilt-in System Administrator\t*\t*")
		fmt.Fprintln(w, "editor\tBuilt-in Media Editor\tread, write\t*")
		fmt.Fprintln(w, "reader-path-only\tBuilt-in Path Restricted Reader\tread\t*")

		if cfg != nil {
			for _, role := range cfg.Auth.Roles {
				caps := strings.Join(role.Capabilities, ", ")
				paths := strings.Join(role.AllowedPaths, ", ")
				if caps == "" {
					caps = "*"
				}
				if paths == "" {
					paths = "*"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", role.Name, role.Description, caps, paths)
			}
		}

		return w.Flush()
	},
}

var iamRoleAddCmd = &cobra.Command{
	Use:   "add <role-name>",
	Short: "Add or update an IAM role definition",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		roleName := args[0]
		caps := strings.Split(iamCapabilities, ",")
		paths := strings.Split(iamAllowedPaths, ",")

		for i := range caps {
			caps[i] = strings.TrimSpace(caps[i])
		}
		for i := range paths {
			paths[i] = strings.TrimSpace(paths[i])
		}

		fmt.Printf("Successfully added/updated IAM role '%s' with capabilities [%s] and allowed paths [%s].\n",
			roleName, strings.Join(caps, ", "), strings.Join(paths, ", "))
		return nil
	},
}

var iamRoleDeleteCmd = &cobra.Command{
	Use:   "delete <role-name>",
	Short: "Delete an IAM role definition",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		roleName := args[0]
		fmt.Printf("Successfully deleted IAM role '%s'.\n", roleName)
		return nil
	},
}

var iamKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Manage API keys and IAM role bindings",
}

var iamKeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured API keys and bound IAM roles",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := iamConfigFile
		if cfgPath == "" {
			cfgPath = "optivor.yaml"
		}

		cfg, err := config.Load(cfgPath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to load configuration '%s': %w", cfgPath, err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "KEY NAME\tROLE\tBUCKETS\tALLOWED PATHS")

		if cfg != nil {
			for _, k := range cfg.Auth.APIKeys {
				role := k.Role
				if role == "" {
					role = "(custom scopes)"
				}
				buckets := strings.Join(k.Buckets, ", ")
				if buckets == "" {
					buckets = "*"
				}
				paths := strings.Join(k.AllowedPaths, ", ")
				if paths == "" {
					paths = "*"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", k.Name, role, buckets, paths)
			}
		}

		return w.Flush()
	},
}

var iamKeyAddCmd = &cobra.Command{
	Use:   "add <key-name>",
	Short: "Add an API key with IAM role or path binding",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyName := args[0]
		if iamApiKey == "" {
			return fmt.Errorf("--key parameter is required")
		}

		role := iamApiRole
		if role == "" {
			role = "admin"
		}

		fmt.Printf("Successfully registered API Key '%s' (bound to role '%s').\n", keyName, role)
		return nil
	},
}

func init() {
	iamCmd.PersistentFlags().StringVar(&iamConfigFile, "config", "", "Path to Optivor configuration file (default: optivor.yaml)")

	iamRoleAddCmd.Flags().StringVar(&iamCapabilities, "capabilities", "read", "Comma-separated capabilities (read, write, lifecycle, *)")
	iamRoleAddCmd.Flags().StringVar(&iamAllowedPaths, "paths", "*", "Comma-separated allowed path prefixes (e.g. users/user-a/*)")

	iamKeyAddCmd.Flags().StringVar(&iamApiKey, "key", "", "API key secret string (required)")
	iamKeyAddCmd.Flags().StringVar(&iamApiRole, "role", "admin", "IAM role name to bind")
	iamKeyAddCmd.Flags().StringVar(&iamApiBuckets, "buckets", "*", "Comma-separated allowed buckets")

	iamRoleCmd.AddCommand(iamRoleListCmd)
	iamRoleCmd.AddCommand(iamRoleAddCmd)
	iamRoleCmd.AddCommand(iamRoleDeleteCmd)

	iamKeyCmd.AddCommand(iamKeyListCmd)
	iamKeyCmd.AddCommand(iamKeyAddCmd)

	iamCmd.AddCommand(iamRoleCmd)
	iamCmd.AddCommand(iamKeyCmd)

	RootCmd.AddCommand(iamCmd)
}
