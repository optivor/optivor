# Frequently Asked Questions (FAQ)

### Why use Optivor instead of cloud image services (Imgix, Cloudinary)?

Cloud image transformation services bill per image or per transform operation, which rapidly scales up costs. Optivor is open-source software you run on your own infrastructure, allowing unlimited image transformations for the cost of standard compute.

### What image formats are supported?

Optivor supports JPEG, PNG, WebP, AVIF, GIF, TIFF, and SVG input formats. Transformed output formats include WebP and AVIF.

### How do I prevent logos or product images from being cropped?

By default, Optivor uses `fit=cover` which center-crops images to fill exact `w` × `h` dimensions. To prevent logo cropping:
- Use `fit=contain` (`?w=500&h=300&fit=contain&format=webp`), which resizes the entire image to fit inside the box preserving aspect ratio.
- Or specify only width (`?w=500&format=webp`), which calculates height proportionally without cropping.
- For full details, see [Image Transformations & Fit Modes](image-transformations.md).

### Why does my downloaded `.webp` file fail to open after using `curl`?

If you run `curl -i ... -o file.webp`, the `-i` flag injects HTTP headers into the file before the binary data, corrupting the image header. Remove `-i` when downloading files:
`curl "http://localhost:8080/image/..." -o file.webp`

### How does caching work?

Optivor includes an in-memory & disk-backed LRU filesystem cache (`internal/cache/fs`). Transformed images are stored locally to serve repeat requests instantly (`X-Optivor-Cache: HIT`). Disk usage is managed automatically by `max_size_mb`.

### Does Optivor support Cloudflare R2 / Backblaze B2 / Google Cloud Storage?

Yes! Any S3-compatible object storage works out of the box with the default S3 driver. External custom drivers can also be managed using `optivor driver install`.
