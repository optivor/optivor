// Optivor Edge CDN Cache Proxy Worker
export default {
  async fetch(request, env, ctx) {
    // Only intercept GET and HEAD requests for CDN caching
    if (request.method !== "GET" && request.method !== "HEAD") {
      const upstreamURL = new URL(new URL(request.url).pathname + new URL(request.url).search, env.OPTIVOR_UPSTREAM_URL);
      return fetch(new Request(upstreamURL.toString(), request));
    }

    const cache = caches.default;

    // 1. Check Cloudflare Edge Cache
    let response = await cache.match(request);
    if (response) {
      const cachedResponse = new Response(response.body, response);
      cachedResponse.headers.set("X-Optivor-Edge-Cache", "HIT");
      return cachedResponse;
    }

    // 2. Fetch from Upstream Optivor Runtime Service
    const url = new URL(request.url);
    const upstreamURL = new URL(url.pathname + url.search, env.OPTIVOR_UPSTREAM_URL);
    const upstreamRequest = new Request(upstreamURL.toString(), request);

    response = await fetch(upstreamRequest);

    // 3. Cache successful image responses at the edge
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
