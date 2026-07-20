# Systemd Deployment Guide

Optivor provides a native systemd deployment adapter to run the runtime as a managed background daemon on Linux servers.

## Quick Installation via Makefile

```bash
make build
sudo make install
```

This installs:
- Binary: `/usr/local/bin/optivor`
- Configuration template: `/etc/optivor/optivor.yaml`
- Systemd unit: `/etc/systemd/system/optivor.service`

## Service Management

Start and enable the service:

```bash
sudo systemctl enable --now optivor
```

Check service status:

```bash
sudo systemctl status optivor
```

View structured JSON logs via journalctl:

```bash
sudo journalctl -u optivor -f -o cat
```

## Uninstallation

```bash
sudo make uninstall
```
