# Edge Integration & CDN Deployment Guide

Optivor is designed to operate seamlessly behind global Content Delivery Networks (CDNs) and Edge Workers such as Cloudflare Workers, AWS CloudFront / Lambda@Edge, and Fastly Compute@Edge.

## Cloudflare Workers Proxy Deployment

Place a Cloudflare Worker in front of your Optivor origin instance to provide global edge caching, DDoS protection, and instant response routing.

### Example Worker Script (`index.js`)

```javascript
export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    const originUrl = `https://optivor.example.com${url.pathname}${url.search}`;

    const cache = caches.default;
    let response = await cache.match(request);

    if (!response) {
      response = await fetch(originUrl, {
        headers: request.headers,
      });

      if (response.status === 200) {
        const headers = new Headers(response.headers);
        headers.set("Cache-Control", "public, max-age=31536000, immutable");
        headers.set("X-Edge-Cache", "MISS");

        response = new Response(response.body, {
          status: response.status,
          headers: headers,
        });

        ctx.waitUntil(cache.put(request, response.clone()));
      }
    } else {
      const headers = new Headers(response.headers);
      headers.set("X-Edge-Cache", "HIT");
      response = new Response(response.body, { headers });
    }

    return response;
  },
};
```

## Bucket Lifecycle Management CLI

Manage retention policies across providers:

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
