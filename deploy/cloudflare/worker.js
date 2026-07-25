export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    const cache = caches.default;

    // Check Cloudflare Edge Cache
    let response = await cache.match(request);
    if (response) {
      return response;
    }

    // Proxy request to upstream Optivor runtime
    const upstreamURL = new URL(url.pathname + url.search, env.OPTIVOR_UPSTREAM_URL);
    const upstreamRequest = new Request(upstreamURL.toString(), request);

    response = await fetch(upstreamRequest);

    // Cache successful image responses at the edge
    if (response.ok && response.headers.get("content-type")?.startsWith("image/")) {
      const responseToCache = new Response(response.body, response);
      responseToCache.headers.set("Cache-Control", "public, max-age=31536000, immutable");
      ctx.waitUntil(cache.put(request, responseToCache.clone()));
      return responseToCache;
    }

    return response;
  }
};
