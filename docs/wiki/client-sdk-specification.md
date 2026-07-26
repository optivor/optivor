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

SDKs must implement both dynamic transformation URLs (`/image/...`) and server-preset URLs (`/preset/{presetName}/...`):

```typescript
interface TransformParams {
  preset?: string;         // Pre-configured server preset (e.g. "avatar", "profile", "thumb")
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

### Preset-First Approach
Optivor encourages preset-first image optimization (`/preset/avatar/user.jpg`). Preset URLs enforce centralized dimensions, quality, and caching rules defined in `optivor.yaml`. All client SDKs support `buildPresetUrl(presetName, key, params)` alongside `buildUrl()`.

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

---

## 5. Ecosystem SDK Usage Guides

### JavaScript / TypeScript (`@optivor/js`)

```bash
npm install @optivor/js
```

```javascript
import { OptivorClient } from '@optivor/js';

const optivor = new OptivorClient({
  baseUrl: 'https://optivor.example.com',
  defaultBucket: 's3-bucket'
});

const url = optivor.buildUrl('users/avatar.jpg', {
  width: 300,
  height: 300,
  fit: 'focal',
  focal: [0.3, 0.7],
  format: 'webp',
  overlay: 'watermark.png',
  gravity: 'bottom_right',
  opacity: 50,
  blur: 5
});
```

### React (`@optivor/react`)

```bash
npm install @optivor/react @optivor/js
```

```jsx
import { OptivorImage } from '@optivor/react';

export default function App() {
  return (
    <OptivorImage
      baseUrl="https://optivor.example.com"
      src="products/shoe.jpg"
      width={400}
      height={300}
      fit="cover"
      format="webp"
      overlay="logo.png"
      gravity="bottom_right"
      opacity={50}
      alt="Sports Shoe"
    />
  );
}
```

### Vue 3 / Nuxt (`@optivor/vue`)

```bash
npm install @optivor/vue @optivor/js
```

```vue
<template>
  <OptivorImage
    base-url="https://optivor.example.com"
    src="products/camera.jpg"
    :width="600"
    :height="400"
    fit="cover"
    format="avif"
    alt="Digital Camera"
  />
</template>

<script setup>
import { OptivorImage } from '@optivor/vue';
</script>
```

### Next.js Custom Loader (`@optivor/next`)

```bash
npm install @optivor/next
```

```tsx
import { Image } from '@optivor/next';

export default function Page() {
  return <Image src="/hero.png" width={1200} height={800} alt="Hero" />;
}
```

### Python (`optivor`)

```bash
pip install optivor
```

```python
from optivor import OptivorClient

optivor = OptivorClient("https://optivor.example.com", "my-bucket")
url = optivor.build_url(
    "photos/landscape.jpg",
    width=800,
    height=600,
    fit="focal",
    focal=(0.4, 0.6),
    format="webp",
    blur=10
)
```

### PHP (`optivor-php`)

```php
use Optivor\OptivorClient;

$optivor = new OptivorClient('https://optivor.example.com', 's3-bucket');
$url = $optivor->url('products/shoes.jpg', [
    'width' => 600,
    'height' => 400,
    'fit' => 'cover',
    'format' => 'webp',
    'overlay' => 'watermark.png',
    'opacity' => 50,
]);
```

