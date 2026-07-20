# Developer Guide: Step-by-Step Driver Implementation

This guide provides step-by-step instructions on implementing a custom storage driver for Optivor.

---

## Step 1: Project Setup & Repository Structure

Create a dedicated repository named `optivor-driver-<provider>` (e.g., `optivor-driver-b2`, `optivor-driver-gcs`).

### Recommended Project Layout (Go example)
```
optivor-driver-b2/
├── .github/
│   └── workflows/
│       └── release.yml
├── cmd/
│   └── optivor-driver-b2/
│       └── main.go
├── internal/
│   └── b2/
│       ├── client.go
│       └── operations.go
├── LICENSE
├── README.md
└── go.mod
```

---

## Step 2: Implement Handshake & CLI Command Parser

Your driver executable must parse the `--optivor-handshake` argument as well as operation commands (`get`, `put`, `head`, `delete`).

### Reference Go Implementation

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type HandshakeInfo struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	OptivorAPI   string   `json:"optivor_api"`
	Capabilities []string `json:"capabilities"`
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--optivor-handshake" {
		info := HandshakeInfo{
			Name:         "b2",
			Version:      "0.1.0",
			OptivorAPI:   "v1",
			Capabilities: []string{"read", "write", "stat", "delete"},
		}
		output, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(output))
		os.Exit(0)
	}

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: optivor-driver-b2 <command> <key>")
		os.Exit(2)
	}

	command := os.Args[1]
	key := os.Args[2]

	switch command {
	case "get":
		handleGet(key)
	case "head":
		handleHead(key)
	default:
		fmt.Fprintf(os.Stderr, "Unsupported command: %s\n", command)
		os.Exit(2)
	}
}

func handleGet(key string) {
	// Read credentials from environment variables
	// endpoint := os.Getenv("OPTIVOR_STORAGE_ENDPOINT")
	// accessKey := os.Getenv("OPTIVOR_STORAGE_ACCESS_KEY")
	
	// Implementation logic streaming object content to os.Stdout...
}

func handleHead(key string) {
	// Output metadata JSON to os.Stdout
}
```

---

## Step 3: Reference Python Implementation

Drivers can be authored in scripting languages like Python (packaged into a single executable using `PyInstaller` or `Nuitka`).

```python
#!/usr/bin/env python3
import sys
import json
import os

HANDSHAKE_RESPONSE = {
    "name": "gcs-custom",
    "version": "0.1.0",
    "optivor_api": "v1",
    "capabilities": ["read", "write", "stat", "delete"]
}

def main():
    if len(sys.argv) > 1 and sys.argv[1] == "--optivor-handshake":
        print(json.dumps(HANDSHAKE_RESPONSE, indent=2))
        sys.exit(0)

    if len(sys.argv) < 3:
        sys.stderr.write("Usage: optivor-driver-gcs <command> <key>\n")
        sys.exit(2)

    command = sys.argv[1]
    key = sys.argv[2]

    if command == "get":
        # Stream content to stdout
        sys.stdout.buffer.write(b"sample image content")
        sys.exit(0)
    elif command == "head":
        meta = {"content_type": "image/jpeg", "content_length": 1024}
        print(json.dumps(meta))
        sys.exit(0)
    else:
        sys.stderr.write(f"Unknown command {command}\n")
        sys.exit(2)

if __name__ == "__main__":
    main()
```

---

## Step 4: Handling Environment Configuration & Credentials

Do not hardcode API keys or bucket names in driver binaries. Always consume standard configuration keys passed down from the Optivor host runtime:

- `OPTIVOR_STORAGE_BUCKET`
- `OPTIVOR_STORAGE_ACCESS_KEY`
- `OPTIVOR_STORAGE_SECRET_KEY`
- `OPTIVOR_STORAGE_ENDPOINT`
- `OPTIVOR_STORAGE_REGION`

Ensure sensible fallbacks or return descriptive error messages on `stderr` with exit code `1` when required credentials are absent.
