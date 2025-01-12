import json

class ResponseParser:
    def __init__(self, response: str) -> None:
        self.response = json.loads(response)

    def parse(self) -> None:
        self.status = self.response['status']
        self.message = self.response['message']
        self.data = self.response['data']
        self.error_code = self.response['error']['code']
        self.error_desc = self.response['error']['description']