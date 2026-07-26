# `@optivor/vue`

> Official Vue 3 & Nuxt component for [Optivor](https://github.com/optivor/optivor) image optimization engine.

[![npm version](https://img.shields.io/npm/v/@optivor/vue.svg)](https://www.npmjs.com/package/@optivor/vue)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/optivor/optivor.svg?style=social)](https://github.com/optivor/optivor)

⭐ **If you find Optivor useful, please consider giving us a [Star on GitHub](https://github.com/optivor/optivor)!**

`<OptivorImage />` is a Vue 3 component for seamless integration with Optivor dynamic image processing servers.

---

## Installation

```bash
npm install @optivor/vue @optivor/js
# or
pnpm add @optivor/vue @optivor/js
```

### Preset-Based Usage (Recommended)

```vue
<template>
  <OptivorImage
    base-url="https://optivor.example.com"
    preset="avatar"
    src="users/john.jpg"
    alt="User Avatar"
  />
</template>

<script setup>
import { OptivorImage } from '@optivor/vue';
</script>
```

---

## Basic Usage

```vue
<template>
  <OptivorImage
    base-url="https://optivor.example.com"
    src="products/camera.jpg"
    :width="600"
    :height="400"
    fit="cover"
    format="webp"
    alt="Digital Camera"
  />
</template>

<script setup>
import { OptivorImage } from '@optivor/vue';
</script>
```

---

## Advanced Options

```vue
<template>
  <OptivorImage
    base-url="https://optivor.example.com"
    src="avatars/user-12.jpg"
    :width="200"
    :height="200"
    fit="focal"
    :focal="[0.5, 0.3]"
    format="avif"
    overlay="verified-badge.png"
    gravity="south_east"
    :opacity="80"
    :grayscale="true"
  />
</template>
```

---

## Watermark Security & Tamper Protection

To ensure watermarks cannot be stripped by removing props or query parameters in public applications:
- **Server Presets (`preset`)**: Pass `preset="watermarked_thumb"` to `<OptivorImage />`. Presets are resolved on the Optivor server without exposing overlay parameters in client URLs.
- **Signed URLs**: When URL signing is enabled on your server, HMAC signatures lock all parameters. Any attempt to modify or remove `overlay` parameters causes a `403 Forbidden` response.

---

## Props Reference

| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `src` | `String` | **(Required)** | Image storage key. |
| `preset` | `String` | `undefined` | Server preset name (e.g. `avatar`, `profile`, `hero`). |
| `baseUrl` / `base-url` | `String` | `'http://localhost:8080'` | Optivor server URL. |
| `width` | `Number` | `undefined` | Target width in pixels. |
| `height` | `Number` | `undefined` | Target height in pixels. |
| `fit` | `String` | `'cover'` | Fit mode (`cover`, `contain`, `fill`, `smart`, `focal`). |
| `format` | `String` | `'webp'` | Target format (`webp`, `avif`, `gif`, `mp4`). |
| `focal` | `Array \| String` | `undefined` | Focal point `[x, y]` coordinates. |
| `overlay` | `String` | `undefined` | Overlay image key. |
| `gravity` | `String` | `undefined` | Overlay gravity position. |
| `opacity` | `Number` | `undefined` | Overlay opacity (0-100). |
| `blur` | `Number` | `undefined` | Gaussian blur radius. |
| `grayscale` | `Boolean` | `false` | Apply monochrome filter. |
| `pixelate` | `Number` | `undefined` | Pixelation block size. |

---

## License

Apache-2.0 © Optivor Team
