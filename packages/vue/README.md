# `@optivor/vue`

> Official Vue 3 & Nuxt component for [Optivor](https://github.com/optivor/optivor) image optimization engine.

[![npm version](https://img.shields.io/npm/v/@optivor/vue.svg)](https://www.npmjs.com/package/@optivor/vue)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

`<OptivorImage />` is a Vue 3 component for seamless integration with Optivor dynamic image processing servers.

---

## Installation

```bash
npm install @optivor/vue @optivor/js
# or
pnpm add @optivor/vue @optivor/js
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

## Props Reference

| Prop | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `src` | `String` | **(Required)** | Image storage key. |
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
