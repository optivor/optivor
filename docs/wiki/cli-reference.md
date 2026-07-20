# CLI Reference

The `optivor` CLI provides commands for scaffolding, deployment, diagnostics, storage driver management, and bucket lifecycle policies.

## Root Command

```bash
optivor [--version] [--help]
```

## Subcommands

### `optivor init`

Scaffolds a new `optivor.yaml` configuration file in the working directory.

```bash
optivor init [--force]
```

### `optivor deploy`

Deploys Optivor using a specified deployment adapter.

```bash
optivor deploy [--adapter systemd] [--config optivor.yaml] [--dry-run]
```

### `optivor doctor`

Performs health and diagnostic checks on system dependencies, configuration, single/multi-bucket connectivity, and libvips runtime.

```bash
optivor doctor [--config optivor.yaml]
```

### `optivor logs`

Tails Optivor service logs via systemd journalctl integration.

```bash
optivor logs [--lines 100] [--follow]
```

### `optivor metrics`

Scrapes and prints runtime metrics from the `/metrics` endpoint.

```bash
optivor metrics [--watch]
```

### `optivor driver`

Manages external storage provider driver binaries.

```bash
# Install driver from local path, github shorthand, or HTTPS release URL
optivor driver install <path-or-url>

# Examples:
optivor driver install /usr/local/bin/optivor-driver-r2
optivor driver install github:optivor/optivor-driver-r2@v1.2.0
optivor driver install https://github.com/optivor/optivor-driver-r2/releases/download/v1.2.0/optivor-driver-r2-linux-amd64

# Driver management
optivor driver list
optivor driver info <name>
optivor driver remove <name>
```

### `optivor bucket lifecycle`

Manages retention policies and expiration lifecycle rules across multi-cloud storage buckets.

```bash
# List active lifecycle rules for a bucket alias
optivor bucket lifecycle list <alias>

# Apply retention policy with custom expiration TTL
optivor bucket lifecycle set <alias> [--ttl-days 30] [--rule-file lifecycle.yaml]

# Delete retention lifecycle rules
optivor bucket lifecycle delete <alias> [--rule-id id] [--all]
```
