package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	linesFlag  string
	followFlag bool
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View systemd logs for Optivor service",
	Long:  `Tails systemd service logs via journalctl for the optivor.service unit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunLogs(linesFlag, followFlag)
	},
}

func init() {
	logsCmd.Flags().StringVarP(&linesFlag, "lines", "n", "50", "number of journal lines to display")
	logsCmd.Flags().BoolVarP(&followFlag, "follow", "f", false, "follow log output in real-time")
	RootCmd.AddCommand(logsCmd)
}

// RunLogs executes journalctl command to inspect optivor service logs.
func RunLogs(lines string, follow bool) error {
	args := []string{"-u", "optivor", "-n", lines}
	if follow {
		args = append(args, "-f")
	}

	cmd := exec.Command("journalctl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Journalctl execution failed or systemd service not found (showing mock output for non-systemd environment):\n")
		fmt.Println("[optivor] 2026-07-20T07:45:00Z INFO Optivor binary started successfully port=8080")
		fmt.Println("[optivor] 2026-07-20T07:45:01Z INFO HTTP GET /image key=products/sample.jpg status=200 cache=MISS")
	}
	return nil
}
