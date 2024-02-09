from datetime import datetime
from typing import Optional, List
from pydantic import BaseModel


class ImageMeta(BaseModel):
    post_id: str
    actor_did: str
    cid: str
    mime_type: str
    created_at: datetime
    data: str


class DetectionResult(BaseModel):
    label: str
    confidence: float
    box: List[float]


class ImageResult(BaseModel):
    meta: ImageMeta
    results: List[DetectionResult]
