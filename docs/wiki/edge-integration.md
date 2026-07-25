# Edge Integration & Cloudflare CDN Deployment Guide

Optivor is designed to operate behind global Content Delivery Networks (CDNs) and Edge Workers such as Cloudflare Workers, AWS CloudFront, and Fastly Compute@Edge.

---

## 1. Cloudflare Edge CDN Integration

The official Cloudflare Edge CDN adapter (`deploy/cloudflare`) intercepts image requests at global Cloudflare edge locations, delivering cached transformed images with sub-10ms latency while reducing load on your origin cluster.

### Architecture Topology

```
                           [ Client Request ]
                                   │
                                   ▼
                   [ Cloudflare Edge CDN Worker ]
                                   │
                 ┌─────────────────┴─────────────────┐
            (Cache Hit)                         (Cache Miss)
                 │                                   │
                 ▼                                   ▼
         [ 200 OK (Edge) ]              [ Origin Optivor Cluster ]
                                      (Kubernetes / Docker / systemd)
```

---

## 2. Upstream Runtime Coupling

### A. Kubernetes (Helm Chart)
Set up Cloudflare Workers in front of a Kubernetes deployment:
1. Expose Optivor via Ingress or LoadBalancer in `deploy/helm/optivor/values.yaml`.
2. Configure `OPTIVOR_UPSTREAM_URL` in `deploy/cloudflare/wrangler.jsonc` pointing to your Ingress URL (e.g. `https://origin-optivor.example.com`).
3. Deploy the Cloudflare Worker using `npx wrangler deploy`.

### B. Docker & Docker Compose
Set up Cloudflare Workers in front of Docker containers:
1. Run Optivor via `docker-compose up -d` listening on port `8080`.
2. Configure a reverse proxy (e.g. Nginx, Caddy, or Cloudflare Tunnel) to route HTTPS traffic to port `8080`.
3. Set `OPTIVOR_UPSTREAM_URL` in `wrangler.jsonc` to your proxy URL.

### C. systemd Service
Set up Cloudflare Workers in front of a systemd daemon:
1. Run Optivor service (`sudo systemctl enable --now optivor`) on `127.0.0.1:8080`.
2. Use Nginx or systemd `cloudflared` to expose the local HTTP port over TLS.
3. Configure `OPTIVOR_UPSTREAM_URL` in `wrangler.jsonc`.

---

## 3. Worker Edge Cache Implementation (`worker.js`)

```javascript
export default {
  async fetch(request, env, ctx) {
    if (request.method !== "GET" && request.method !== "HEAD") {
      const upstreamURL = new URL(new URL(request.url).pathname + new URL(request.url).search, env.OPTIVOR_UPSTREAM_URL);
      return fetch(new Request(upstreamURL.toString(), request));
    }

    const cache = caches.default;
    let response = await cache.match(request);
    if (response) {
      const cachedResponse = new Response(response.body, response);
      cachedResponse.headers.set("X-Optivor-Edge-Cache", "HIT");
      return cachedResponse;
    }

    const url = new URL(request.url);
    const upstreamURL = new URL(url.pathname + url.search, env.OPTIVOR_UPSTREAM_URL);
    response = await fetch(new Request(upstreamURL.toString(), request));

    if (response.ok && response.headers.get("content-type")?.startsWith("image/")) {
      const responseToCache = new Response(response.body, response);
      responseToCache.headers.set("Cache-Control", "public, max-age=31536000, s-maxage=31536000, immutable");
      responseToCache.headers.set("X-Optivor-Edge-Cache", "MISS");
      ctx.waitUntil(cache.put(request, responseToCache.clone()));
      return responseToCache;
    }

    return response;
  }
};
```

---

## 4. Bucket Lifecycle Management CLI

Manage retention policies across cloud storage providers:

```bash
# List current rules
optivor bucket lifecycle list primary-images

# Apply 30-day expiration policy
optivor bucket lifecycle set primary-images --ttl-days 30

# Apply policy from YAML definition
optivor bucket lifecycle set primary-images --rule-file lifecycle.yaml

# Clear lifecycle rules
optivor bucket lifecycle delete primary-images --all
```
