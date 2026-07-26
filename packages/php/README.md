# `optivor-php`

> Official PHP & Laravel SDK for [Optivor](https://github.com/optivor/optivor) image optimization engine.

[![Packagist Version](https://img.shields.io/packagist/v/optivor/optivor-php.svg)](https://packagist.org/packages/optivor/optivor-php)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/optivor/optivor.svg?style=social)](https://github.com/optivor/optivor)

⭐ **If you find Optivor useful, please consider giving us a [Star on GitHub](https://github.com/optivor/optivor)!**

`optivor-php` provides URL generation and parameter formatting for PHP 7.4+ applications, Laravel Blade directives, and Symfony integrations.

---

## Installation

```bash
composer require optivor/optivor-php
```

---

## Usage

```php
use Optivor\OptivorClient;

$optivor = new OptivorClient('https://optivor.example.com', 's3-bucket');

```php
// Preset-based URL generation (Recommended)
$presetUrl = $optivor->presetUrl('avatar', 'users/john.jpg');
// => https://optivor.example.com/preset/avatar/s3-bucket/users/john.jpg

// Basic optimized WebP URL
$url = $optivor->url('products/shoes.jpg', [
    'width' => 600,
    'height' => 400,
    'fit' => 'cover',
    'format' => 'webp',
]);
// => https://optivor.example.com/image/s3-bucket/products/shoes.jpg?w=600&h=400&fit=cover&format=webp
```

---

## Advanced Options

```php
$url = $optivor->url('portraits/model.jpg', [
    'width' => 800,
    'height' => 800,
    'fit' => 'focal',
    'focal' => [0.4, 0.7],
    'format' => 'avif',
    'overlay' => 'watermark.png',
    'gravity' => 'bottom_right',
    'opacity' => 50,
    'blur' => 10,
    'grayscale' => true,
]);
```

---

## Watermark Security & Tamper Protection

To prevent end-users from stripping `overlay=` or watermark query parameters from public URLs:

1. **HMAC URL Signing**: Generate signed URLs using `securityKey`. HMAC signatures lock the URL path and parameters (`overlay=`). Any attempt to alter or remove the watermark parameter invalidates the signature, returning `403 Forbidden`.
2. **Server-Side Presets (`presetUrl`)**: Use `$optivor->presetUrl('watermarked_thumb', $key)` to process images using server-side presets defined in `optivor.yaml`. The overlay rule stays on the server, leaving no parameters for users to tamper with.

---

## License

Apache-2.0 © Optivor Team
