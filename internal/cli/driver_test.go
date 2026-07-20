package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createFakeDriver(t *testing.T, dir string, name string, version string, apiVersion string) string {
	t.Helper()
	fakePath := filepath.Join(dir, "optivor-driver-"+name)
	scriptContent := "#!/bin/sh\n" +
		"if [ \"$1\" = \"--optivor-handshake\" ]; then\n" +
		"  echo '{\"name\":\"" + name + "\",\"version\":\"" + version + "\",\"optivor_api\":\"" + apiVersion + "\"}'\n" +
		"  exit 0\n" +
		"fi\n"
	err := os.WriteFile(fakePath, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatalf("failed to create fake driver script: %v", err)
	}
	return fakePath
}

func TestDriverSubcommands(t *testing.T) {
	tempDir := t.TempDir()
	registryFilePath = filepath.Join(tempDir, "drivers.json")
	defer func() { registryFilePath = "" }()

	fakeDriverPath := createFakeDriver(t, tempDir, "test-r2", "1.0.0", "v1")

	t.Run("driver install", func(t *testing.T) {
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"driver", "install", fakeDriverPath})

		err := RootCmd.Execute()
		if err != nil {
			t.Fatalf("expected install success, got error: %v", err)
		}
	})

	t.Run("driver list", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		RootCmd.SetArgs([]string{"driver", "list"})
		err := RootCmd.Execute()

		w.Close()
		os.Stdout = oldStdout

		if err != nil {
			t.Fatalf("expected list success, got error: %v", err)
		}

		var outBuf bytes.Buffer
		outBuf.ReadFrom(r)
		outStr := outBuf.String()
		if !strings.Contains(outStr, "test-r2") || !strings.Contains(outStr, "1.0.0") {
			t.Errorf("expected output to contain 'test-r2' and '1.0.0', got: %s", outStr)
		}
	})

	t.Run("driver info", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		RootCmd.SetArgs([]string{"driver", "info", "test-r2"})
		err := RootCmd.Execute()

		w.Close()
		os.Stdout = oldStdout

		if err != nil {
			t.Fatalf("expected info success, got error: %v", err)
		}

		var outBuf bytes.Buffer
		outBuf.ReadFrom(r)
		outStr := outBuf.String()
		if !strings.Contains(outStr, "Name:        test-r2") {
			t.Errorf("expected info output to contain 'Name:        test-r2', got: %s", outStr)
		}
	})

	t.Run("driver remove", func(t *testing.T) {
		RootCmd.SetArgs([]string{"driver", "remove", "test-r2"})
		err := RootCmd.Execute()
		if err != nil {
			t.Fatalf("expected remove success, got error: %v", err)
		}

		// Verify empty list
		reg, err := loadRegistry()
		if err != nil {
			t.Fatalf("failed to load registry: %v", err)
		}
		if len(reg.Drivers) != 0 {
			t.Errorf("expected 0 drivers after removal, got %d", len(reg.Drivers))
		}
	})
}
