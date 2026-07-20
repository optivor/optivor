# CLI Reference

The `optivor` CLI provides commands for scaffolding, deployment, diagnostics, and storage driver management.

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

Performs health and diagnostic checks on system dependencies, configuration, S3 connectivity, and libvips runtime.

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
optivor driver install <path>
optivor driver list
optivor driver info <name>
optivor driver remove <name>
```
