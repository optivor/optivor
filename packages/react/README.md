# `@optivor/react`

> Official React SDK & `<OptivorImage />` component for [Optivor](https://github.com/optivor/optivor) image optimization engine.

[![npm version](https://img.shields.io/npm/v/@optivor/react.svg)](https://www.npmjs.com/package/@optivor/react)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/optivor/optivor.svg?style=social)](https://github.com/optivor/optivor)

⭐ **If you find Optivor useful, please consider giving us a [Star on GitHub](https://github.com/optivor/optivor)!**

`<OptivorImage />` provides automatic responsive `srcset` generation, lazy loading, and zero-config Optivor URL formatting for React applications.

---

## Installation

```bash
npm install @optivor/react @optivor/js
# or
pnpm add @optivor/react @optivor/js
```

### Preset-Based Usage (Recommended)

```jsx
// Use server-side preset (e.g. /preset/avatar/users/alex.jpg)
<OptivorImage
  baseUrl="https://optivor.example.com"
  preset="avatar"
  src="users/alex.jpg"
  alt="User Avatar"
/>
```

---

## Usage

```jsx
import { OptivorImage } from '@optivor/react';

export default function UserCard() {
  return (
    <OptivorImage
      baseUrl="https://optivor.example.com"
      src="users/alex.jpg"
      width={300}
      height={300}
      fit="cover"
      format="webp"
      alt="User Profile"
      className="rounded-full shadow-lg"
    />
  );
}
```

### Advanced Effects & Watermarking

```jsx
<OptivorImage
  baseUrl="https://optivor.example.com"
  src="products/shoe.jpg"
  width={800}
  height={600}
  fit="focal"
  focal={[0.4, 0.6]}
  format="avif"
  overlay="logo.png"
  gravity="bottom_right"
  opacity={50}
  blur={5}
  alt="Sports Shoe"
/>
```

---

## Watermark Security & Anti-Tamper Protection

To ensure watermarks cannot be stripped by removing props or query parameters in public applications:
- **Server Presets (`preset`)**: Pass `preset="watermarked_thumb"` to `<OptivorImage />`. Presets are resolved on the Optivor server without exposing overlay parameters in client URLs.
- **Signed URLs**: When URL signing is enabled on your server, HMAC signatures lock all parameters. Any attempt to modify or remove `overlay` parameters causes a `403 Forbidden` response.

---

## Props Reference

| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `src` | `string` | **(Required)** | Image key or path in Optivor storage. |
| `preset` | `string` | `undefined` | Optivor server preset name (e.g. `avatar`, `profile`, `thumb`). |
| `baseUrl` | `string` | `'http://localhost:8080'` | Base URL of your Optivor server. |
| `width` | `number` | `undefined` | Display & target width in pixels. Generates `srcset` automatically. |
| `height` | `number` | `undefined` | Target height in pixels. |
| `fit` | `'cover' \| 'contain' \| 'fill' \| 'smart' \| 'focal'` | `'cover'` | Resizing strategy. |
| `format` | `'webp' \| 'avif' \| 'gif' \| 'mp4'` | `'webp'` | Output format. |
| `focal` | `[number, number]` | `undefined` | Focal point coordinates `[x, y]` (0.0 to 1.0). |
| `overlay` | `string` | `undefined` | Filigran overlay image key. |
| `gravity` | `string` | `'center'` | Position of overlay. |
| `opacity` | `number` | `100` | Overlay opacity (0 to 100). |
| `blur` | `number` | `undefined` | Blur radius. |
| `grayscale` | `boolean` | `false` | Apply black-and-white filter. |
| `pixelate` | `number` | `undefined` | Pixelation block size. |

---

## License

Apache-2.0 © Optivor Team
