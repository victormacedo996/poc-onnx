from pydantic import BaseModel

class HealthResponse(BaseModel):
    message: str


class PredictResponse(BaseModel):
    prediction: str