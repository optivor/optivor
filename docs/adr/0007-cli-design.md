# ADR-0007: CLI Architecture and Command Structure

## Status

Accepted

## Context

As Optivor transitions to V0.3, developers require a unified Command Line Interface (CLI) to initialize local project environments (`optivor init`), inspect deployment configurations, and orchestrate deployment adapters (`optivor deploy`).

Per **ADR-0002** (High-Level Architecture) and **ADR-0003** (Extensibility Model), the CLI must observe strict layer boundaries:
- The CLI layer must **never** process image data directly or interface with storage backends directly.
- The CLI orchestrates external adapters (e.g. systemd) via out-of-process subprocess execution or system invocation.
- Subpackages under `cmd/optivor` or `internal/cli` must remain decoupled from `internal/pipeline` and `internal/cache`.

## Decision

1. **CLI Binary Entrypoint**:
   The `optivor` executable serves both as the runtime server (when run without subcommands or with `--config`) and as the CLI runner (when subcommands like `init` or `deploy` are specified).

2. **CLI Framework**:
   We select **Cobra** (`github.com/spf13/cobra`) as the CLI routing and flag handling framework. Cobra is lightweight, standard across Go infrastructure tools, and strictly aligns with ADR-0001 (Minimal Dependency Policy).

3. **Command Surface (V0.3 Scope)**:
   - `optivor init`: Generates a standard `optivor.yaml` configuration file and `.gitignore` safety recommendations in the current directory.
   - `optivor deploy`: Dispatches deployment actions to configured deployment adapters (defaulting to `--adapter systemd`). Supports `--dry-run` to preview actions without making system modifications.
   - `optivor --version`: Prints the current version string.

4. **Adapter Execution Boundary**:
   The CLI invokes adapter commands out-of-process (e.g. executing systemd installation scripts or Makefile helpers) without directly importing adapter packages, enforcing ADR-0003 decoupled extensibility.

## Consequences

- End-users obtain a consistent developer experience for bootstrapping and deploying Optivor.
- Architectural layer discipline remains strictly enforced, preventing CLI logic from polluting server runtime performance or pipeline code.
