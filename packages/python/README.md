# `optivor` (Python SDK)

> Official Python SDK for [Optivor](https://github.com/optivor/optivor) image optimization engine (Django, FastAPI, Flask).

[![PyPI version](https://img.shields.io/pypi/v/optivor.svg)](https://pypi.org/project/optivor/)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/optivor/optivor.svg?style=social)](https://github.com/optivor/optivor)

⭐ **If you find Optivor useful, please consider giving us a [Star on GitHub](https://github.com/optivor/optivor)!**

The `optivor` Python package provides helper methods for building signed and dynamic image transformation URLs.

---

## Installation

```bash
pip install optivor
```

---

## Basic Usage

```python
from optivor import OptivorClient

optivor = OptivorClient(
    base_url="https://optivor.example.com",
    default_bucket="my-bucket"
)

```python
# Preset-based image optimization (Recommended)
preset_url = optivor.build_preset_url("avatar", "users/john.jpg")
# => https://optivor.example.com/preset/avatar/my-bucket/users/john.jpg

# Generate WebP image URL
url = optivor.build_url(
    "photos/landscape.jpg",
    width=800,
    height=600,
    fit="cover",
    format="webp"
)
# => https://optivor.example.com/image/my-bucket/photos/landscape.jpg?w=800&h=600&fit=cover&format=webp
```

---

## Advanced Options

```python
url = optivor.build_url(
    "users/avatar.jpg",
    width=400,
    height=400,
    fit="focal",
    focal=(0.3, 0.7),
    format="avif",
    overlay="logo.png",
    gravity="bottom_right",
    opacity=60,
    blur=5.0,
    grayscale=True,
    pixelate=4
)
```

---

## Watermark Security & Tamper Protection

To prevent end-users from stripping `overlay=` or watermark query parameters from public URLs:

1. **HMAC URL Signing**: Generate signed URLs using `security_key`. The HMAC-SHA256 signature locks both the path and query parameters (`overlay=`). Modifying or removing the overlay parameter invalidates the signature, returning `403 Forbidden`.
2. **Server-Side Presets (`build_preset_url`)**: Use `optivor.build_preset_url('watermarked_thumb', key)` to process images using server-side presets defined in `optivor.yaml`. The overlay rule stays on the server, leaving no parameters for users to tamper with.

---

## Django Integration Example

```python
# settings.py
OPTIVOR_URL = "https://optivor.example.com"
OPTIVOR_BUCKET = "prod-images"

# templatetags/optivor_tags.py
from django import template
from optivor import OptivorClient
from django.conf import settings

register = template.Library()
optivor = OptivorClient(settings.OPTIVOR_URL, settings.OPTIVOR_BUCKET)

@register.simple_tag
def optivor_url(key, **kwargs):
    return optivor.build_url(key, **kwargs)
```

---

## License

Apache-2.0 © Optivor Team
