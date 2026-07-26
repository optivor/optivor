"""
Optivor Django Template Filter Integration
Usage in Django Templates:
    {% load optivor_tags %}
    <img src="{% optivor_url 'avatar' photo.key %}" />
"""

try:
    from django import template
    from django.conf import settings
    from optivor import OptivorClient

    register = template.Library()

    def _get_client():
        base_url = getattr(settings, 'OPTIVOR_BASE_URL', 'http://localhost:8080')
        security_key = getattr(settings, 'OPTIVOR_SECURITY_KEY', None)
        default_bucket = getattr(settings, 'OPTIVOR_DEFAULT_BUCKET', 'default')
        return OptivorClient(base_url=base_url, security_key=security_key, default_bucket=default_bucket)

    @register.simple_tag
    def optivor_url(preset_or_key: str, key: str = None, bucket: str = None, **kwargs) -> str:
        """Generates an Optivor image URL within Django templates."""
        client = _get_client()
        if key is None:
            # Usage: {% optivor_url 'image.jpg' w=800 h=600 %}
            return client.build_url(key=preset_or_key, bucket=bucket, **kwargs)
        else:
            # Usage: {% optivor_url 'avatar' 'photo.jpg' %}
            return client.preset_url(preset=preset_or_key, key=key, bucket=bucket)

except ImportError:
    # Django is optional dependency
    pass
