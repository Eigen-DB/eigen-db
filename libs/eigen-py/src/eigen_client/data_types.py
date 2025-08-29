class Embedding:
    def __init__(self, id: int, data: list[float], metadata: dict[str, str] = {}) -> None:
        self.id = id
        self.data = data
        self.metadata = metadata

    def to_dict(self) -> dict:
        return {
            "id": self.id,
            "data": self.data,
            "metadata": self.metadata
        }
    
    def __repr__(self) -> str:
        return f"Embedding(id={self.id}, data={self.data}, metadata={self.metadata})"
    
class Document():
    def __init__(self, id: int, data: str, metadata: dict[str, str] = {}) -> None:
        self.id = id
        self.data = data
        self.metadata = metadata

    def __repr__(self) -> str:
        return f"Document(id={self.id} data={self.data}, metadata={self.metadata})"
