from urllib.parse import urlencode, rstrip

class OptivorClient:
    def __init__(self, base_url: str = "http://localhost:8080", default_bucket: str = ""):
        self.base_url = base_url.rstrip("/")
        self.default_bucket = default_bucket

    def build_url(self, key: str, **params) -> str:
        clean_key = key.lstrip("/")
        if self.default_bucket and "/" not in clean_key:
            full_path = f"{self.default_bucket}/{clean_key}"
        else:
            full_path = clean_key

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

        q_str = urlencode(query)
        return f"{self.base_url}/image/{full_path}{'?' + q_str if q_str else ''}"
