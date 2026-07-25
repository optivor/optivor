# Optivor Cloudflare Edge CDN Cache Integration Guide

This adapter provides a high-performance **Cloudflare Worker Edge CDN Cache** in front of your Optivor runtime cluster (Kubernetes, Docker, or systemd).

---

## 1. Overview & Architecture

The Cloudflare Worker operates as an **Edge CDN Cache Proxy**:
- Intercepts image requests at global Cloudflare edge locations (`caches.default`).
- Delivers cached images directly to end users with sub-10ms latency.
- Adds `X-Optivor-Edge-Cache: HIT` / `MISS` diagnostic headers.
- Proxies cache misses to your upstream Optivor runtime cluster and caches the optimized response.

```
                           [ User Request ]
                                  │
                                  ▼
                   [ Cloudflare Edge CDN Worker ]
                                  │
                  ┌───────────────┴───────────────┐
             (Cache Hit)                     (Cache Miss)
                  │                               │
                  ▼                               ▼
          [ 200 OK (Edge) ]          [ Upstream Optivor Runtime ]
                                     (Kubernetes / Docker / systemd)
```

---

## 2. Upstream Environment Coupling

### A. Kubernetes (Helm Deployment)

If Optivor is deployed on Kubernetes using the official Helm chart (`deploy/helm/optivor`):

1. **Expose Optivor Service via Ingress / LoadBalancer**:
   Ensure `values.yaml` has ingress enabled with a public domain or AWS ALB hostname:
   ```yaml
   ingress:
     enabled: true
     className: "nginx"
     hosts:
       - host: optivor-origin.yourdomain.com
         paths:
           - path: /
             pathType: Prefix
   ```

2. **Configure Cloudflare Worker Upstream**:
   In `deploy/cloudflare/wrangler.jsonc`:
   ```json
   {
     "name": "optivor-edge-cdn",
     "main": "worker.js",
     "compatibility_date": "2026-07-25",
     "vars": {
       "OPTIVOR_UPSTREAM_URL": "https://optivor-origin.yourdomain.com"
     }
   }
   ```

3. **Alternative (Zero-Trust Cloudflare Tunnel)**:
   Deploy a `cloudflared` deployment inside Kubernetes routing traffic directly to `http://optivor.default.svc.cluster.local:8080` without exposing public ports.

---

### B. Docker & Docker Compose

If Optivor is running in Docker:

1. **Docker Compose Setup (`docker-compose.yml`)**:
   Expose port `8080` behind a reverse proxy (e.g. Nginx, Caddy, or `cloudflared` container):
   ```yaml
   services:
     optivor:
       image: optivor/optivor:latest
       ports:
         - "8080:8080"
       volumes:
         - ./optivor.yaml:/etc/optivor/optivor.yaml:ro
   ```

2. **Configure Upstream Domain**:
   Point `OPTIVOR_UPSTREAM_URL` in `wrangler.jsonc` to your server's domain or public IP:
   ```json
   "vars": {
     "OPTIVOR_UPSTREAM_URL": "https://optivor-origin.yourdomain.com"
   }
   ```

---

### C. systemd Daemon (Linux Bare-Metal / VM)

If Optivor is running as a systemd service:

1. **Service Verification**:
   Ensure the service is active on `http://127.0.0.1:8080`:
   ```bash
   sudo systemctl status optivor
   ```

2. **Reverse Proxy (Nginx / Caddy)**:
   Pass incoming requests from port `443` (SSL) to `http://127.0.0.1:8080`:
   ```nginx
   server {
       server_name optivor-origin.yourdomain.com;
       location / {
           proxy_pass http://127.0.0.1:8080;
           proxy_set_header Host $host;
       }
   }
   ```

3. **Configure Worker**:
   Set `OPTIVOR_UPSTREAM_URL` in `wrangler.jsonc` to `https://optivor-origin.yourdomain.com`.

---

## 3. Quick Start & Deployment

1. Install Wrangler CLI:
   ```bash
   npm install -g wrangler
   ```

2. Set your upstream URL in `wrangler.jsonc`.

3. Deploy to Cloudflare Workers:
   ```bash
   npx wrangler deploy
   ```
