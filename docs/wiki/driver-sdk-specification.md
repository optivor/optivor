# Storage Driver SDK & IPC Protocol Specification

This specification defines the communication protocol between Optivor runtime and external storage driver executables.

## Handshake Protocol (`--optivor-handshake`)

Drivers MUST support the `--optivor-handshake` flag and return a JSON object on standard output:

```json
{
  "name": "r2",
  "version": "1.0.0",
  "optivor_api": "v1"
}
```

## Binary IPC Interface

Optivor executes the driver binary with parameters passed via standard command-line flags or environment variables:

```bash
optivor-driver-r2 get --key "images/sample.jpg" --endpoint "..." --bucket "..."
```

The driver MUST output the raw object byte stream to stdout (fd 1) and exit with code 0 on success. On failure, it MUST output error messages to stderr (fd 2) and exit with a non-zero exit code.
