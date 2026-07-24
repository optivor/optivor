# `@optivor/next`

Zero-config Next.js `next/image` component and loader integration for Optivor.

## Installation

```bash
npm install @optivor/next
# or
pnpm add @optivor/next
```

## Zero-Config Component Usage (Recommended)

Set environment variable:

```env
NEXT_PUBLIC_OPTIVOR_URL=https://optivor.example.com
```

Use `<Image />` directly in your Next.js App Router or Pages Router app with zero `next.config.js` edits:

```tsx
import { Image } from '@optivor/next';

export default function Hero() {
  return (
    <Image
      src="/hero.png"
      width={1200}
      height={800}
      alt="Hero Image"
    />
  );
}
```

## Custom Loader Usage

If you prefer using standard `<Image />` tags with Next.js loader configuration:

```js
// optivor-loader.js
const { optivorLoader } = require('@optivor/next');
module.exports = optivorLoader;
```

```js
// next.config.js
module.exports = {
  images: {
    loader: 'custom',
    loaderFile: './optivor-loader.js',
  },
};
```
