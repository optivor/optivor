# Optivor Cloudflare Edge CDN Cache Adapter

This adapter provides a zero-overhead **Cloudflare Worker Edge CDN Cache** in front of your Optivor runtime cluster.

## Overview

It operates strictly as a high-performance **Edge CDN Cache**:
- Intercepts image requests at global Cloudflare edge locations (`caches.default`).
- Serves cached transformed images directly from the edge with sub-10ms response times.
- Forwards cache misses to your upstream Optivor runtime and caches the optimized result.
- Passes non-GET requests directly to upstream.

## Quick Start

1. Install Wrangler CLI (if not already installed):
   ```bash
   npm install -g wrangler
   ```

2. Configure upstream URL in `wrangler.jsonc`:
   ```json
   "vars": {
     "OPTIVOR_UPSTREAM_URL": "https://optivor.yourdomain.com"
   }
   ```

3. Deploy to Cloudflare Workers:
   ```bash
   npx wrangler deploy
   ```
