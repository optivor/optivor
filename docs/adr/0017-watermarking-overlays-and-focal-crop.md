# ADR-0017: Watermarking, Image Overlays, Focal Point Cropping & Dynamic Filters

- **Status**: Accepted
- **Deciders**: Optivor Core Architecture Team
- **Date**: 2026-07-26

---

## Context and Problem Statement

Modern web applications and e-commerce platforms require advanced media transformation capabilities directly within the image optimization pipeline:
1. **Dynamic Watermarking & Overlays**: Brand protection and asset attribution require overlaying logos or watermarks with custom positioning (`gravity`), transparency (`opacity`), and relative scaling (`overlay_scale`).
2. **Manual Focal Point Cropping**: While smart entropy cropping (`fit=smart`) automatically detects high-energy regions, user-generated content (portraits, product focal points) often specifies explicit focal coordinates (`focal=0.3,0.7`) to keep key subjects centered.
3. **Animated Media Handling**: Converting animated GIFs to high-efficiency animated WebP or video formats reduces bandwidth significantly for micro-animations.
4. **Image Filters & Effects**: Applications need dynamic image enhancements such as Gaussian blur (`blur=N`), grayscale conversion (`grayscale=true`), and pixelation (`pixelate=N`) for preview placeholders and UI background effects.

---

## Decision Outcome

We decide to expand Optivor's `pipeline.TransformParams` and libvips processing pipeline with native support for watermarking/overlays, manual focal point cropping, animated media handling, and dynamic image filters:

### 1. Watermarking & Overlays (`CompositeImage`)
- **Query Parameters**: `overlay=<key_or_url>`, `gravity=<position>`, `opacity=<0-100>`, `overlay_scale=<1-100>`.
- **Supported Gravities**: `center`, `north`, `south`, `east`, `west`, `north_east`, `north_west`, `south_east`, `south_west`, `bottom_right`, `bottom_left`, `top_right`, `top_left`.
- **Implementation**: libvips compositing (`Composite2`) with linear alpha blend and position calculation based on calculated target bounding boxes.

### 2. Manual Focal Point Cropping (`FocalCrop`)
- **Query Parameter**: `focal=<x>,<y>` (normalized coordinates between `0.0` and `1.0`, e.g., `focal=0.3,0.7`).
- **Fit Mode**: `fit=focal` (or automatic fallback when `focal` parameter is supplied).
- **Implementation**: Calculates target cropping window relative to focal point coordinates `(x * image_width, y * image_height)` while respecting image bounds.

### 3. Animated Media Conversion
- **Supported Formats**: Animated WebP, GIF, MP4/WebM micro-conversions.
- **Implementation**: Load multi-frame images (`n=-1`) into govips and encode with frame rate metadata preservation.

### 4. Dynamic Image Filters
- **Blur**: `blur=<1-100>` applies libvips Gaussian blur.
- **Grayscale**: `grayscale=true` converts color space to monochrome interpretation (`InterpretationBW`).
- **Pixelate**: `pixelate=<1-50>` applies downsampling and block upsampling for privacy blurring and pixel art placeholders.

---

## Concurrency & Performance Impact

- All composite operations leverage in-memory libvips buffers to prevent disk I/O bottlenecks.
- Parameter validation guards against extreme values (e.g. maximum blur radius capped to 100, maximum opacity bounded to 0-100%).

---

## Consequences

### Positive
- Unified API for watermarking, focal cropping, animated media, and visual filter effects.
- Zero external service dependencies; processing remains high-performance and in-process via libvips.
- Full compatibility with existing URL signatures, caching layers, and bucket routing engines.

### Negative
- Compositing and multi-frame processing slightly increase memory allocation per request, mitigated by resource bounds.
