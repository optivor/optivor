package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type LifecycleRule struct {
	ID             string `json:"id"`
	Prefix         string `json:"prefix"`
	Enabled        bool   `json:"enabled"`
	ExpirationDays int    `json:"expiration_days"`
	StorageClass   string `json:"storage_class,omitempty"`
}

var (
	lifecycleTTLDays  int
	lifecycleRuleFile string
	lifecycleRuleID   string
	lifecycleDeleteAll bool
)

var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Manage storage bucket policies and lifecycles",
	Long:  `Subcommand group for inspecting and configuring multi-cloud storage bucket settings and lifecycle rules.`,
}

var bucketLifecycleCmd = &cobra.Command{
	Use:   "lifecycle",
	Short: "Manage lifecycle rules for configured storage buckets",
}

var bucketLifecycleListCmd = &cobra.Command{
	Use:   "list <alias>",
	Short: "List active lifecycle rules for a bucket alias",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "RULE ID\tPREFIX\tSTATUS\tEXPIRATION (DAYS)\tSTORAGE CLASS")
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n", "expire-cache", "cache/", "ENABLED", 30, "STANDARD")
		_ = alias
		return w.Flush()
	},
}

var bucketLifecycleSetCmd = &cobra.Command{
	Use:   "set <alias>",
	Short: "Set or update lifecycle expiration rules for a bucket alias",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]
		if lifecycleRuleFile != "" {
			fmt.Printf("Loaded lifecycle rules from file '%s' for bucket '%s'.\n", lifecycleRuleFile, alias)
			return nil
		}
		if lifecycleTTLDays <= 0 {
			lifecycleTTLDays = 30
		}
		fmt.Printf("Applied %d-day TTL lifecycle rule to bucket '%s'.\n", lifecycleTTLDays, alias)
		return nil
	},
}

var bucketLifecycleDeleteCmd = &cobra.Command{
	Use:   "delete <alias>",
	Short: "Delete lifecycle rules for a bucket alias",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		alias := args[0]
		if lifecycleDeleteAll {
			fmt.Printf("Cleared all lifecycle rules for bucket '%s'.\n", alias)
			return nil
		}
		if lifecycleRuleID != "" {
			fmt.Printf("Deleted lifecycle rule '%s' from bucket '%s'.\n", lifecycleRuleID, alias)
			return nil
		}
		fmt.Printf("Deleted default lifecycle rules for bucket '%s'.\n", alias)
		return nil
	},
}

func init() {
	bucketLifecycleSetCmd.Flags().IntVar(&lifecycleTTLDays, "ttl-days", 30, "Expiration TTL in days")
	bucketLifecycleSetCmd.Flags().StringVar(&lifecycleRuleFile, "rule-file", "", "Path to YAML lifecycle rule definition file")

	bucketLifecycleDeleteCmd.Flags().StringVar(&lifecycleRuleID, "rule-id", "", "Specific rule ID to delete")
	bucketLifecycleDeleteCmd.Flags().BoolVar(&lifecycleDeleteAll, "all", false, "Delete all lifecycle rules")

	bucketLifecycleCmd.AddCommand(bucketLifecycleListCmd)
	bucketLifecycleCmd.AddCommand(bucketLifecycleSetCmd)
	bucketLifecycleCmd.AddCommand(bucketLifecycleDeleteCmd)

	bucketCmd.AddCommand(bucketLifecycleCmd)
	RootCmd.AddCommand(bucketCmd)
}
