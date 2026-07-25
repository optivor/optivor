# Interactive CLI Wizard Guide (`optivor init --interactive`)

The Optivor CLI provides a guided interactive wizard to simplify configuration generation for new environments and local setups.

---

## Usage

To launch the interactive setup wizard:

```bash
optivor init --interactive
# or short alias
optivor init -i
```

---

## Features

1. **Guided Config Generation**:
   Prompts step-by-step for server port, log levels, storage endpoints, access credentials, and cache options.

2. **Automated Validation**:
   Verifies that entered values conform to expected configuration schemas before writing `optivor.yaml`.

3. **Overwrite Safety**:
   Prevents accidental overwrites of existing `optivor.yaml` files unless explicitly passed `--force`:

```bash
optivor init --interactive --force
```

---

## Output Example

```text
=== Optivor Interactive Configuration Wizard ===
Scaffolding default S3 storage and server configuration...
Successfully created optivor.yaml via interactive wizard!

Next steps:
  1. Verify configuration with 'optivor doctor'
  2. Start Optivor server with 'optivor'
```
