import hmac
import hashlib
from urllib.parse import urlencode

class OptivorClient:
    def __init__(self, base_url: str = "http://localhost:8080", security_key: str = None, default_bucket: str = ""):
        self.base_url = base_url.rstrip("/")
        self.security_key = security_key
        self.default_bucket = default_bucket

    def build_preset_url(self, preset_name: str, key: str, **params) -> str:
        clean_key = key.lstrip("/")
        full_path = f"{self.default_bucket}/{clean_key}" if self.default_bucket and "/" not in clean_key else clean_key
        q_str = self._build_query(**params)
        raw_url = f"{self.base_url}/preset/{preset_name}/{full_path}{'?' + q_str if q_str else ''}"
        if self.security_key:
            sig = hmac.new(self.security_key.encode('utf-8'), raw_url.encode('utf-8'), hashlib.sha256).hexdigest()[:16]
            join_char = "&" if "?" in raw_url else "?"
            return f"{raw_url}{join_char}s={sig}"
        return raw_url

    def preset_url(self, preset_name: str, key: str, **params) -> str:
        return self.build_preset_url(preset_name, key, **params)

    def build_url(self, key: str, **params) -> str:
        if "preset" in params and params["preset"]:
            preset_name = params.pop("preset")
            return self.build_preset_url(preset_name, key, **params)

        clean_key = key.lstrip("/")
        full_path = f"{self.default_bucket}/{clean_key}" if self.default_bucket and "/" not in clean_key else clean_key
        q_str = self._build_query(**params)
        raw_url = f"{self.base_url}/image/{full_path}{'?' + q_str if q_str else ''}"
        if self.security_key:
            sig = hmac.new(self.security_key.encode('utf-8'), raw_url.encode('utf-8'), hashlib.sha256).hexdigest()[:16]
            join_char = "&" if "?" in raw_url else "?"
            return f"{raw_url}{join_char}s={sig}"
        return raw_url

    def _build_query(self, **params) -> str:
        query = {}
        if "width" in params or "w" in params:
            query["w"] = params.get("width") or params.get("w")
        if "height" in params or "h" in params:
            query["h"] = params.get("height") or params.get("h")
        if "fit" in params:
            query["fit"] = params["fit"]
        if "format" in params:
            query["format"] = params["format"]
        if "focal" in params:
            focal = params["focal"]
            query["focal"] = ",".join(map(str, focal)) if isinstance(focal, (list, tuple)) else str(focal)
        if "overlay" in params:
            query["overlay"] = params["overlay"]
        if "gravity" in params:
            query["gravity"] = params["gravity"]
        if "opacity" in params:
            query["opacity"] = params["opacity"]
        if "blur" in params:
            query["blur"] = params["blur"]
        if params.get("grayscale"):
            query["grayscale"] = "true"
        if "pixelate" in params:
            query["pixelate"] = params["pixelate"]

        return urlencode(query)
