from requests import put, delete, post, get
import json

from eigen_client import API_VERSION
from eigen_client.index import Index
from eigen_client.response import ResponseParser
from eigen_client.utils import validate_embedding_model

class Client:
    def __init__(self, 
        url: str,
        api_key: str
    ) -> None:
        self.url = url
        self.base_url = f'{url}/api/{API_VERSION}'
        self.api_key = api_key
        self.headers = {
            "X-Eigen-API-Key": self.api_key,
            "Content-Type": "application/json"
        }
        self._test_auth()

    def _test_auth(self) -> bool:
        '''
        Test the provided API key against the EigenDB instance.
        Returns:
            True if the API key is valid, False otherwise.
        '''
        res = get(url=self.base_url + '/test-auth', headers=self.headers)
        try:
            parser = ResponseParser(res)
            parser.parse()
        except Exception as e:
            print(f"Authentication failed: {e}")
            return False
        
        print("Authentication successful:", parser.message)
        return True
    
    def create_index_from_model(self, index_name: str, model_name: str, model_provider: str, model_provider_api_key: str = None) -> Index:
        """
        Create a new index using a supported embedding model. 
        If the model is not supported, use create_index instead with the desired dimensions and metric.

        Args:
            index_name (str): The name of the index to create.
            model_name (str): The name of the model to use for embeddings.

        Returns:
            Index: An instance of the Index class.
        """
        model = validate_embedding_model(model_name, model_provider)
        self.create_index(
            index_name=index_name,
            dimensions=model['dimensions'],
            metric=model['metric']
        )
        return Index(
            api_key=self.api_key,
            url=self.url,
            index_name=index_name,
            model_provider=model_provider,
            model_name=model_name,
            model_provider_api_key=model_provider_api_key
        )

    def create_index(self, index_name: str, dimensions: int, metric: str) -> Index:
        """
        Create a new index with the specified dimensions and metric.

        Args:
            index_name (str): The name of the index to create.
            dimensions (int): The number of dimensions for the embeddings.
            metric (str): The distance metric to use (e.g., "cosine", "euclidean").

        Returns:
            Index: An instance of the Index class.
        """
        url = f"{self.base_url}/indexes/{index_name}/create"
        body = {
            "dimensions": dimensions,
            "metric": metric
        }
        response = put(url, headers=self.headers, data=json.dumps(body))
        parser = ResponseParser(response)
        parser.parse()
        return Index(
            api_key=self.api_key,
            url=self.url,
            index_name=index_name,
            model_provider="none",
        )
    
    def delete_index(self, index_name: str) -> None:
        """
        Delete an existing index.

        Args:
            index_name (str): The name of the index to delete.
        """
        url = f"{self.base_url}/indexes/{index_name}/delete"
        response = delete(url, headers=self.headers)
        parser = ResponseParser(response)
        parser.parse()

    def list_indexes(self) -> list[str]:
        """
        List all existing indexes.

        Returns:
            list[str]: A list of index names.
        """
        url = f"{self.base_url}/indexes/list"
        response = get(url, headers=self.headers)
        parser = ResponseParser(response)
        parser.parse()
        return parser.data["indexes"]
    
    def get_index_stats(self, index_name: str) -> dict:
        """
        Get statistics for a specific index.

        Args:
            index_name (str): The name of the index.
        Returns:
            dict: A dictionary containing index statistics.
        """
        url = f"{self.base_url}/indexes/{index_name}/stats"
        response = get(url, headers=self.headers)
        parser = ResponseParser(response)
        parser.parse()
        return parser.data
    
    def __repr__(self) -> str:
        return f"Client(url={self.base_url})"
