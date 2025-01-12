from pydantic import BaseModel
from typing import Optional

class ErrorDetails(BaseModel):
    code: str
    description: str

class Response(BaseModel):
    status: int
    message: str
    error: Optional[ErrorDetails] = None