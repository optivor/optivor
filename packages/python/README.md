# `optivor` (Python SDK)

> Official Python SDK for [Optivor](https://github.com/optivor/optivor) image optimization engine (Django, FastAPI, Flask).

[![PyPI version](https://img.shields.io/pypi/v/optivor.svg)](https://pypi.org/project/optivor/)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

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
