from fastapi import FastAPI

from app.api.resume import router as resume_router
from app.api.health import router as health_router


def create_app() -> FastAPI:
    """Create and configure the FastAPI application.

    Inputs: none.
    Outputs: a FastAPI app with the health router registered.
    """

    app = FastAPI(
        title="Interview Coach AI Service",
        version="0.1.0",
        docs_url="/docs",
        redoc_url="/redoc",
    )
    app.include_router(health_router)
    app.include_router(resume_router)
    return app


app = create_app()
