# Frequently Asked Questions (FAQ)

### Why use Optivor instead of cloud image services (Imgix, Cloudinary)?

Cloud image transformation services bill per image or per transform operation, which rapidly scales up costs. Optivor is open-source software you run on your own infrastructure, allowing unlimited image transformations for the cost of standard compute.

### What image formats are supported?

Optivor supports JPEG, PNG, WebP, AVIF, GIF, TIFF, and SVG input formats. Transformed output formats include WebP and AVIF.

### How does caching work?

Optivor includes an in-memory & disk-backed LRU filesystem cache (`internal/cache/fs`). Transformed images are stored locally to serve repeat requests instantly (`X-Optivor-Cache: HIT`). Disk usage is managed automatically by `max_size_mb`.

### Does Optivor support Cloudflare R2 / Backblaze B2 / Google Cloud Storage?

Yes! Any S3-compatible object storage works out of the box with the default S3 driver. External custom drivers can also be managed using `optivor driver install`.
