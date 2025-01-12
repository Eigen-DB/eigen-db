from requests import get, post, put
import json

from response import ResponseParser

class Client:
    def __init__(self, url: str, api_key: str) -> None:
        self.url = url
        self.api_key = api_key

    def insert_vector(self, id: int, embedding: list[float]) -> None:
        res = put(
            url=self.url,
            headers={
                'X-Eigen-API-Key': self.api_key
            },
            data=json.dumps({
                "vector": {
                    "id": id,
                    "embedding": embedding
                }
            })
        )
        parser = ResponseParser(res.content)
        parser.parse()

        if res.status_code != 200:
            raise Exception(f"{parser.error_code} - {parser.error_desc}")
        
        