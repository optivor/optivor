# Image Transformations & Fit Modes

Optivor provides high-performance, real-time image transformation and format conversion. This guide details transformation parameters, resizing modes (`fit`), preset configurations, and CLI testing best practices.

## Transformation Parameters

| Query Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `w` | integer | original | Target width in pixels (e.g. `w=500`). |
| `h` | integer | original | Target height in pixels (e.g. `h=300`). |
| `fit` | string | `cover` | Resizing mode: `cover`, `contain`, `fill`. |
| `format` | string | original | Target output format: `webp` or `avif`. |

---

## Fit Modes (`fit`)

Understanding the `fit` parameter is essential for achieving the desired visual presentation (e.g. avoiding unexpected logo cropping).

### 1. `fit=cover` (Default)
- **Behavior:** Scales the image to fill the exact dimensions (`w` × `h`), performing a **center crop** if the aspect ratio of the source image does not match the target dimensions.
- **Use Case:** User profile avatars, card background images, thumbnail grids where consistent box dimensions are required.
- **Example:**
  ```bash
  curl "http://localhost:8080/image/r2-bucket/sample.jpg?w=500&h=300&fit=cover&format=webp" -o cover.webp
  ```

### 2. `fit=contain`
- **Behavior:** Resizes the image to fit entirely within the specified dimensions (`w` × `h`) while strictly **preserving the original aspect ratio**. No cropping occurs. Any remaining space is padded with a background (transparent for WebP/AVIF/PNG).
- **Use Case:** Company logos, product photos, icons, and brand graphics where full visibility without cropping is required.
- **Example:**
  ```bash
  curl "http://localhost:8080/image/r2-bucket/logo.png?w=500&h=300&fit=contain&format=webp" -o contain.webp
  ```

### 3. `fit=fill`
- **Behavior:** Scales the image to exact `w` × `h` dimensions, stretching or squeezing the image if necessary without preserving aspect ratio.
- **Use Case:** Specific UI containers where exact aspect ratio stretching is intended.

### 4. Proportional Resizing (Single Dimension)
If you omit either `w` or `h` (e.g. provide only `w=500`), Optivor automatically calculates the missing dimension based on the original image's aspect ratio, preventing any cropping:
```bash
curl "http://localhost:8080/image/r2-bucket/logo.png?w=500&format=webp" -o proportional.webp
```

---

## Preset Endpoints (`/preset/...`)

Presets allow you to define pre-configured transformation sets in `optivor.yaml` to standardize image variants across your frontend applications:

```yaml
presets:
  avatar:
    w: 150
    h: 150
    f: "webp"
    fit: "cover"
  banner:
    w: 1200
    h: 400
    f: "webp"
    fit: "cover"
```

### Invoking Presets
Request preset transformations using the `/preset/{name}/{bucket-alias}/{key}` route:

```bash
curl "http://localhost:8080/preset/avatar/r2-bucket/users/alex.jpg" -o avatar.webp
```

---

## CLI & cURL Testing Pro Tips

> [!IMPORTANT]
> **Avoid `-i` when downloading image binary files:**
> When using `curl` with `--output` or `-o`, do **NOT** include the `-i` (include HTTP headers) flag.
> The `-i` flag writes ASCII HTTP headers (e.g., `HTTP/1.1 200 OK`) directly to the top of the binary file, corrupting the image file structure and causing image viewers to reject it.

### Correct Download Syntax:
```bash
# Correct: Saves clean WebP binary file
curl "http://localhost:8080/image/r2-bucket/photo.jpg?w=400&format=webp" -o output.webp

# Inspecting HTTP Headers only:
curl -i "http://localhost:8080/image/r2-bucket/photo.jpg?w=400&format=webp"
```

### Verifying Generated Image Metadata:
Use the Linux `file` utility to verify that the generated image is a valid WebP/AVIF file:
```bash
file output.webp
# Expected output: output.webp: RIFF (little-endian) data, Web/P image
```
