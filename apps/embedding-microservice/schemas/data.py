from pydantic import BaseModel

class Data(BaseModel):
    id: int
    data: str