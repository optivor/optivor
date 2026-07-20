package cli

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var (
	metricsEndpointFlag string
	watchFlag           bool
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Scrape and display Prometheus metrics from Optivor runtime",
	Long:  `Connects to the /metrics endpoint of a running Optivor server and displays raw or summary metrics.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunMetrics(metricsEndpointFlag, watchFlag)
	},
}

func init() {
	metricsCmd.Flags().StringVarP(&metricsEndpointFlag, "endpoint", "e", "http://localhost:8080/metrics", "metrics endpoint URL")
	metricsCmd.Flags().BoolVarP(&watchFlag, "watch", "w", false, "refresh metrics display periodically")
	RootCmd.AddCommand(metricsCmd)
}

// RunMetrics fetches and displays metrics from target endpoint.
func RunMetrics(endpoint string, watch bool) error {
	fetch := func() error {
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(endpoint)
		if err != nil {
			return fmt.Errorf("failed to fetch metrics from %s: %w", endpoint, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("metrics endpoint returned status %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read metrics body: %w", err)
		}

		fmt.Printf("--- Optivor Metrics Output (%s) ---\n", endpoint)
		fmt.Println(string(body))
		return nil
	}

	if err := fetch(); err != nil {
		return err
	}

	return nil
}
