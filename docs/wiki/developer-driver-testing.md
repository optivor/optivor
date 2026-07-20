# Developer Guide: Driver Testing & Verification

This guide explains how to test and validate your custom Optivor storage driver before submitting it to the community registry.

---

## 1. Handshake Verification Test

Verify that your binary returns valid JSON and exits with code `0` when called with `--optivor-handshake`.

```bash
./optivor-driver-b2 --optivor-handshake
```

### Expected Output
```json
{
  "name": "b2",
  "version": "0.1.0",
  "optivor_api": "v1",
  "capabilities": ["read", "write", "stat", "delete"]
}
```

### Verification Script
```bash
# Verify exit code 0
./optivor-driver-b2 --optivor-handshake > /dev/null
if [ $? -ne 0 ]; then
  echo "FAIL: Handshake exit code must be 0"
  exit 1
fi

# Verify JSON schema validation
./optivor-driver-b2 --optivor-handshake | jq .optivor_api
```

---

## 2. Local Registration Test

Use the Optivor CLI to register your binary locally and confirm driver discovery.

```bash
# 1. Install driver binary
optivor driver install ./optivor-driver-b2

# 2. List installed drivers
optivor driver list

# 3. View driver metadata info
optivor driver info b2
```

---

## 3. End-to-End Operation Tests

Set test environment variables and test individual operations (`head`, `get`, `put`, `delete`).

```bash
export OPTIVOR_STORAGE_BUCKET="test-bucket"
export OPTIVOR_STORAGE_ACCESS_KEY="test-key"
export OPTIVOR_STORAGE_SECRET_KEY="test-secret"

# Test HEAD operation (Metadata check)
./optivor-driver-b2 head sample.jpg

# Test GET operation (Payload stream)
./optivor-driver-b2 get sample.jpg > output.jpg

# Test PUT operation (Payload write)
cat input.jpg | ./optivor-driver-b2 put sample.jpg

# Test DELETE operation
./optivor-driver-b2 delete sample.jpg
```

---

## 4. Error Handling Checklist

Ensure your driver handles common edge cases appropriately:

| Scenario | Expected Behavior | Exit Code |
|---|---|---|
| Non-existent key (`get`) | Print error to `stderr` | `1` |
| Invalid credentials | Print authentication error to `stderr` | `1` |
| Missing required CLI argument | Print usage guide to `stderr` | `2` |
| Network timeout | Retry internally or print error to `stderr` | `1` |
