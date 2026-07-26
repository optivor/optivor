package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestIAMCmd(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "optivor.yaml")

	cfgContent := `
storage:
  driver: local
auth:
  roles:
    - name: "custom-editor"
      description: "Custom Editor Role"
      capabilities: ["read", "write"]
      allowed_paths: ["media/*"]
  api_keys:
    - key: "secret-key-1"
      name: "User A Key"
      role: "custom-editor"
      buckets: ["tenant-bucket"]
      allowed_paths: ["media/*"]
`
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0600); err != nil {
		t.Fatalf("failed to create temp config: %v", err)
	}

	t.Run("iam role list", func(t *testing.T) {
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"iam", "role", "list", "--config", cfgPath})
		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error running iam role list: %v", err)
		}
	})

	t.Run("iam role add", func(t *testing.T) {
		RootCmd.SetArgs([]string{"iam", "role", "add", "tenant-role", "--capabilities", "read", "--paths", "users/user-a/*"})
		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error running iam role add: %v", err)
		}
	})

	t.Run("iam role delete", func(t *testing.T) {
		RootCmd.SetArgs([]string{"iam", "role", "delete", "tenant-role"})
		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error running iam role delete: %v", err)
		}
	})

	t.Run("iam key list", func(t *testing.T) {
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"iam", "key", "list", "--config", cfgPath})
		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error running iam key list: %v", err)
		}
	})

	t.Run("iam key add", func(t *testing.T) {
		RootCmd.SetArgs([]string{"iam", "key", "add", "new-user-key", "--key", "secret123", "--role", "custom-editor"})
		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("unexpected error running iam key add: %v", err)
		}
	})
}
