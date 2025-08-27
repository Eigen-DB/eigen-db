from requests import Response
import json

class ResponseParser:
    '''
    Parses the response from an EigenDB instance.
    Args:
        response: The Response object from the requests library received from a request to the EigenDB API.
    '''

    def __init__(self, response: Response) -> None:
        if response.status_code != 200:
            error_message = json.loads(response.content.decode('utf-8'))
            raise Exception(f"An error occured: {response.status_code} {response.reason}\nError details:\n{json.dumps(error_message, indent=4)}")
        
        self.response = json.loads(response.content.decode('utf-8'))
        self.status: int = -1
        self.message: str = ""
        self.data: dict = {}
        self.error_code: str = ""
        self.error_desc: str = ""

    def parse(self) -> None:
        '''
        Parses the response and extracts relevant information.
        '''
        self.status = self.response['status'] if 'status' in self.response else -1
        self.message = self.response['message'] if 'message' in self.response else ""
        self.data = self.response['data'] if 'data' in self.response else {}
        self.error_code = self.response['error']['code'] if 'error' in self.response else ""
        self.error_desc = self.response['error']['description'] if 'error' in self.response else ""