# Official Client SDK Specification

This specification establishes the standard contract for official Optivor framework SDKs across JavaScript, React, Vue, PHP, and Python.

---

## 1. Official Client Packages

| Ecosystem / Platform | Package Identifier | Purpose |
| :--- | :--- | :--- |
| Core JavaScript / TS | `@optivor/js` | URL builder, signature generator, client options helper. |
| React / Next.js | `@optivor/react` | React `<OptivorImage />` component with responsive `srcset`. |
| Vue / Nuxt | `@optivor/vue` | Vue 3 `<OptivorImage>` component & plugin integration. |
| PHP / Laravel | `optivor-php` | PHP Client & Laravel service provider / Blade directive. |
| Python / Django | `optivor-python` | Python SDK & Django / FastAPI URL building helpers. |

---

## 2. Standardized Configuration & Options

All SDKs accept a uniform set of configuration properties:

```typescript
interface OptivorClientOptions {
  baseUrl: string;         // e.g. "https://optivor.example.com"
  securityKey?: string;    // Secret key for HMAC signed URLs
  defaultBucket?: string;  // Default storage bucket alias
  defaultFormat?: "webp" | "avif" | "gif" | "mp4";
}
```

---

## 3. URL Construction Contract

SDKs must implement a deterministic `buildUrl(key, params)` method:

```typescript
interface TransformParams {
  width?: number;
  height?: number;
  fit?: 'cover' | 'contain' | 'fill' | 'smart' | 'focal';
  format?: 'webp' | 'avif' | 'gif' | 'mp4';
  focal?: [number, number];
  overlay?: string;
  gravity?: string;
  opacity?: number;
  blur?: number;
  grayscale?: boolean;
  pixelate?: number;
}
```

### Signature Generation
When `securityKey` is configured:
1. Append expiration timestamp `expires=<unix_timestamp>` (if provided).
2. Generate SHA-256 HMAC signature over the raw path and query string.
3. Append `s=<signature>` parameter.

---

## 4. UI Component Specification (`@optivor/react` / `@optivor/vue`)

Components must render responsive `<img>` or `<picture>` elements with automatic `srcset` generation:

```jsx
<OptivorImage
  src="products/shoe.jpg"
  width={400}
  height={300}
  fit="cover"
  format="webp"
  blur={5}
  alt="Sports shoe"
/>
```

Generated HTML:
```html
<img
  src="https://optivor.example.com/image/products/shoe.jpg?w=400&h=300&fit=cover&format=webp&blur=5"
  srcset="https://optivor.example.com/image/products/shoe.jpg?w=400&... 1x, https://optivor.example.com/image/products/shoe.jpg?w=800&... 2x"
  width="400"
  height="300"
  alt="Sports shoe"
  loading="lazy"
/>
```
