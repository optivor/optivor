"""
Optivor FastAPI Dependency Injector Integration
Usage in FastAPI:
    from fastapi import FastAPI, Depends
    from optivor.fastapi import get_optivor_client, OptivorClient

    app = FastAPI()

    @app.get("/avatar")
    def get_avatar(optivor: OptivorClient = Depends(get_optivor_client)):
        return {"url": optivor.preset_url("avatar", "user.jpg")}
"""

import os
from optivor import OptivorClient

def get_optivor_client() -> OptivorClient:
    """FastAPI dependency provider for OptivorClient instance."""
    base_url = os.getenv("OPTIVOR_BASE_URL", "http://localhost:8080")
    security_key = os.getenv("OPTIVOR_SECURITY_KEY")
    default_bucket = os.getenv("OPTIVOR_DEFAULT_BUCKET", "default")
    return OptivorClient(base_url=base_url, security_key=security_key, default_bucket=default_bucket)
