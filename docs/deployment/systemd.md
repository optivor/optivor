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

---

## Cloudflare Edge CDN Integration

To pair your systemd daemon (`127.0.0.1:8080`) with Cloudflare Workers as a global Edge CDN:

1. Reverse proxy traffic from your domain `https://optivor-origin.yourdomain.com` to `http://127.0.0.1:8080` via Nginx/Caddy or `cloudflared`.
2. Configure `deploy/cloudflare/wrangler.jsonc`:
   ```json
   "vars": {
     "OPTIVOR_UPSTREAM_URL": "https://optivor-origin.yourdomain.com"
   }
   ```
3. Deploy the worker using `npx wrangler deploy` to enable global edge caching.
