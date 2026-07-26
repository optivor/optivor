# Watermarking, Overlays & Visual Effects Guide

This guide details Optivor's dynamic watermarking, overlay compositing, manual focal point cropping, animated media conversion, and visual filter options.

---

## 1. Dynamic Watermarking & Overlays

Optivor allows dynamically compositing overlay images (watermarks, logos, brand badges) over optimized base images.

### Query Parameters

| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `overlay` | string | - | Path or key of the overlay image (e.g. `watermark.png`). |
| `gravity` | string | `center` | Positioning anchor: `north_west`, `north`, `north_east`, `west`, `center`, `east`, `south_west`, `south`, `south_east` (or alias `bottom_right`, `top_left`, etc.). |
| `opacity` | float | `100` | Overlay transparency level (0-100%). |
| `overlay_scale` | float | - | Scale overlay relative to base image width (1-100%). |

### Example Usage
```bash
curl "http://localhost:8080/image/default/photo.jpg?w=800&overlay=watermark.png&gravity=bottom_right&opacity=50&overlay_scale=20&format=webp" -o watermarked.webp
```

---

## 2. Manual Focal Point Cropping (`focal`)

Focal point cropping provides direct coordinate-based cropping (`focal=x,y`) to complement automated smart entropy cropping (`fit=smart`).

### Query Parameter
- `focal=<x>,<y>`: Normalized horizontal and vertical coordinates between `0.0` and `1.0`.

### Example
```bash
# Crop target 200x200 centered around focal point at 30% width, 70% height
curl "http://localhost:8080/image/default/portrait.jpg?w=200&h=200&fit=focal&focal=0.3,0.7&format=webp" -o focal.webp
```

---

## 3. Animated Media Conversion

Optivor automatically converts multi-frame GIFs and animated images to compressed animated formats or video micro-animations.

### Supported Formats
- `format=webp`: Multi-frame animated WebP output.
- `format=gif`: Optimized GIF output.
- `format=mp4`: Micro-animation video stream with `video/mp4` MIME output.

### Example
```bash
curl "http://localhost:8080/image/default/banner.gif?w=400&format=mp4" -o animation.mp4
```

---

## 4. Visual Image Filters

Optivor includes built-in visual filters for placeholder creation, UI background blurs, and privacy effects.

### Query Parameters

| Parameter | Type | Range | Description |
| :--- | :--- | :--- | :--- |
| `blur` | float | `0.1 - 100` | Gaussian blur radius in pixels (e.g. `blur=10`). |
| `grayscale` | boolean | `true/false` | Converts color space to monochrome (`grayscale=true`). |
| `pixelate` | integer | `2 - 50` | Block downscaling factor for pixel art / privacy blurring (e.g. `pixelate=8`). |

### Example Usage
```bash
curl "http://localhost:8080/image/default/hero.jpg?w=600&blur=15&grayscale=true&pixelate=4&format=webp" -o filtered.webp
```

---

## 5. Watermark Security & Anti-Tamper Protection

When serving protected or copyrighted media, preventing clients from stripping `overlay=` or watermark parameters from the URL is critical. Optivor provides two zero-trust architectural mechanisms:

### A. HMAC URL Signatures (`securityKey`)
When URL signing is enabled on Optivor (`securityKey` in config or SDK):
- The server generates an HMAC-SHA256 signature calculated over the entire URL path **and** query string parameters (including `overlay=...`, `opacity=...`, `gravity=...`).
- If an end-user or bad actor attempts to strip `overlay=watermark.png` from the URL, the signature check fails immediately with **`403 Forbidden`**.

### B. Server-Side Presets (`/preset/{presetName}/{key}`)
For maximum security without exposing transform query parameters:
- Define a preset in `optivor.yaml` (e.g. `watermarked_preview` with `overlay: "logo.png"`).
- Frontend applications request `/preset/watermarked_preview/photos/item.jpg`.
- Because the overlay configuration resides strictly on the Optivor engine, there are no `overlay` query parameters present in the URL string to strip.

