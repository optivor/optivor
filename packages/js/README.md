# `@optivor/js`

> Official JavaScript & TypeScript SDK for [Optivor](https://github.com/optivor/optivor) image optimization engine.

[![npm version](https://img.shields.io/npm/v/@optivor/js.svg)](https://www.npmjs.com/package/@optivor/js)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

`@optivor/js` is a lightweight, zero-dependency client library for constructing Optivor image transformation URLs, handling signed URL generation, and managing multi-bucket image routing.

---

## Installation

```bash
npm install @optivor/js
# or
pnpm add @optivor/js
# or
yarn add @optivor/js
```

---

## Basic Usage

```javascript
import { OptivorClient } from '@optivor/js';

const optivor = new OptivorClient({
  baseUrl: 'https://optivor.example.com',
  defaultBucket: 's3-bucket'
});

// Generate optimized WebP URL
const imageUrl = optivor.buildUrl('users/avatar.jpg', {
  width: 300,
  height: 300,
  fit: 'cover',
  format: 'webp'
});
// => https://optivor.example.com/image/s3-bucket/users/avatar.jpg?w=300&h=300&fit=cover&format=webp
```

---

## Advanced Transformation Options

`@optivor/js` supports all Optivor V1.2 transformation parameters:

```javascript
const bannerUrl = optivor.buildUrl('banners/hero.png', {
  width: 1200,
  height: 600,
  fit: 'focal',
  focal: [0.3, 0.7],          // Focal point X: 30%, Y: 70%
  format: 'avif',
  overlay: 'watermark.png',   // Filigran / logo overlay key
  gravity: 'bottom_right',    // Overlay position
  opacity: 50,                // Transparency (0-100%)
  blur: 10,                   // Gaussian blur radius
  grayscale: true,            // B&W conversion
  pixelate: 5                 // Block pixelation
});
```

---

## Options & Parameters Reference

### `OptivorClientOptions`

| Option | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `baseUrl` | `string` | `'http://localhost:8080'` | Base URL of your Optivor server. |
| `defaultBucket` | `string` | `''` | Optional default storage bucket alias. |

### `TransformParams`

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `width` / `w` | `number` | Target image width in pixels. |
| `height` / `h` | `number` | Target image height in pixels. |
| `fit` | `'cover' \| 'contain' \| 'fill' \| 'smart' \| 'focal'` | Resizing & cropping strategy. |
| `format` | `'webp' \| 'avif' \| 'gif' \| 'mp4'` | Output image format. |
| `focal` | `[number, number]` | Normalized focal point coordinates `[x, y]` (0.0 to 1.0). |
| `overlay` | `string` | Overlay image key / logo path. |
| `gravity` | `string` | Position: `center`, `north_west`, `bottom_right`, etc. |
| `opacity` | `number` | Overlay transparency (0 to 100). |
| `blur` | `number` | Gaussian blur radius. |
| `grayscale` | `boolean` | Convert image to monochrome. |
| `pixelate` | `number` | Downscale/upscale block size for pixel art. |

---

## License

Apache-2.0 © Optivor Team
