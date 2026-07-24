# `@optivor/next`

Zero-config Next.js `next/image` custom loader integration for Optivor.

## Installation

```bash
npm install @optivor/next
# or
pnpm add @optivor/next
```

## Usage

In your `next.config.js`:

```js
module.exports = {
  images: {
    loader: 'custom',
    loaderFile: './node_modules/@optivor/next/index.js',
  },
};
```

Set environment variable:

```env
NEXT_PUBLIC_OPTIVOR_URL=https://optivor.example.com
```

Now all standard `<Image />` tags will automatically route through your Optivor instance for high-performance WebP/AVIF transformation and caching!
