from fastapi import APIRouter
from pydantic import BaseModel

router = APIRouter()


class HealthResponse(BaseModel):
    """HealthResponse describes the service health payload."""

    status: str


@router.get("/health", response_model=HealthResponse)
async def health() -> HealthResponse:
    """Return the AI service health response.

    Inputs: none.
    Outputs: a JSON object with `status=ok` when the service is healthy.
    """

    return HealthResponse(status="ok")
