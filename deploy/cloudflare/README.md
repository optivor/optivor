# Optivor Cloudflare Edge Deployment Adapter

This deployment adapter provisions a **Cloudflare Worker Edge Proxy** in front of the core Optivor runtime container, following the deployment adapter architecture specified in **ADR-0002** and **ADR-0012**.

## Architecture

Per **ADR-0002**:
> "Deployment Adapters explicitly document whether they deploy the full runtime or a proxy in front of it."

This adapter deploys an **Edge Proxy**. Cloudflare Workers handle edge caching, URL signed token validation, rate-limiting, and geo-routing, while passing uncached image optimization requests to the upstream Optivor runtime service.

```
[ Client Request ] ---> [ Cloudflare Worker (Edge Cache & Auth) ] ---> [ Optivor Container Runtime ]
```

## Configuration

### `wrangler.jsonc`

```json
{
  "$schema": "node_modules/wrangler/config-schema.json",
  "name": "optivor-edge-proxy",
  "main": "src/worker.js",
  "compatibility_date": "2026-07-25",
  "vars": {
    "OPTIVOR_UPSTREAM_URL": "https://optivor.yourdomain.com"
  }
}
```

## Quick Start

```bash
cd deploy/cloudflare
npm install
npx wrangler deploy
```
