from requests import get, put, delete
from typing import Literal
import json
import ollama
import openai
import tiktoken

from eigen_client.response import ResponseParser
from eigen_client.data_types import Embedding, Document
from eigen_client.supported_models import SUPPORTED_MODEL_PROVIDERS, SUPPORTED_MODELS

#logging.set_verbosity_error()

class Index:
    '''
    EigenDB API client for managing embeddings.
    Args:
        url: The URL of the EigenDB instance.
        api_key: The API key for the EigenDB instance.
        model_provider: The embedding model provider to use for vectorization.
        model_name: The name of the embedding model to use for vectorization (disregard if model_provider = "none").
        model_provider_api_key: The API key for the embedding model provider (if required).
    '''

    def __init__(self, 
        url: str,
        api_key: str,
        model_provider: Literal["openai", "ollama", "none"] = "openai",
        model_name: str = "text-embedding-3-small",
        model_provider_api_key: str = None,
    ) -> None:
        # checking model params
        model_names = [model['name'] for model in SUPPORTED_MODELS if model_provider in model['supported_providers']]
        if model_provider not in SUPPORTED_MODEL_PROVIDERS:
            raise ValueError(f"Invalid model provider: {model_provider}. Supported providers are: {', '.join(SUPPORTED_MODEL_PROVIDERS)}.")
        if model_name not in model_names:
            raise ValueError(f"Invalid model name: {model_name} for provider {model_provider}. Supported models for this provider are: {', '.join(model_names)}.")

        self.url = url
        self.api_key = api_key
        self.model = [model for model in SUPPORTED_MODELS if model['name'] == model_name][0]
        self.model_provider = model_provider
        self.model_provider_api_key = model_provider_api_key
        
        self._test_auth()
        if self.model_provider == "openai":
            self.openai_client = openai.Client(api_key=self.model_provider_api_key)
        elif self.model_provider == "ollama":
            self.ollama_client = ollama

    def _test_auth(self) -> bool:
        '''
        Test the provided API key against the EigenDB instance.
        Returns:
            True if the API key is valid, False otherwise.
        '''
        res = get(
            url=self.url + '/test-auth',
            headers={
                'X-Eigen-API-Key': self.api_key
            }
        )
        try:
            parser = ResponseParser(res)
            parser.parse()
        except Exception as e:
            print(f"Authentication failed: {e}")
            return False
        
        print("Authentication successful:", parser.message)
        return True
    
    def _vectorize_docs(self, docs: list[Document]) -> list[Embedding]:
        '''
        Vectorize a given list of documents using the specified model provider and model name.
        Args:
            docs: A list of Document objects to be vectorized.
        Returns:
            A list of Embedding objects containing the vectorized representations of the input documents.
        '''
        output_embeddings: list[Embedding] = []
        if self.model_provider == "ollama":
            response = self.ollama_client.embed(
                model=self.model['name'],
                input=[doc.data for doc in docs],
                truncate=True
            )
            for i in range(len(response.embeddings)):
                output_embeddings.append(
                    Embedding(
                        id=docs[i].id,
                        data=response.embeddings[i],
                        metadata=docs[i].metadata
                    )
                )
        elif self.model_provider == "openai":
            enc = tiktoken.get_encoding("cl100k_base") # should work for all openai embedding models we support at the moment: https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
            token_limit = self.model['metadata']['token_limit']
            for doc in docs:
                encoded_doc = enc.encode(doc.data)
                num_tokens = len(encoded_doc)
                if num_tokens > token_limit:
                    print(f"Warning: Document with ID {doc.id} exceeds the token limit of {token_limit} tokens for model {self.model['name']}. It has {num_tokens} tokens and will be truncated.")
                    truncated_encoded_doc = encoded_doc[:token_limit]
                    doc.data = enc.decode(truncated_encoded_doc)

            response = self.openai_client.embeddings.create(
                model=self.model['name'],
                input=[doc.data for doc in docs],
            )
            for i in range(len(response.data)):
                output_embeddings.append(
                    Embedding(
                        id=docs[i].id,
                        data=response.data[i].embedding,
                        metadata=docs[i].metadata
                    )
                )
        elif self.model_provider == "none":
            raise Exception("No model provider specified. Please set a model provider to vectorize documents.")
        else:
            raise Exception(f"Invalid model provider: {self.model_provider}. Supported providers are: {', '.join(SUPPORTED_MODEL_PROVIDERS)}.")

        return output_embeddings

    def insert(self, embeddings: list[Embedding]) -> None:
        '''
        Inserts a set of embeddings into the EigenDB instance.
        Args:
            embeddings: A list of Embedding objects to be inserted.
        '''
        res = put(
            url=self.url + '/embeddings/insert',
            headers={
                'X-Eigen-API-Key': self.api_key
            },
            data=json.dumps({
                "embeddings": [e.to_dict() for e in embeddings]
            })
        )
        parser = ResponseParser(res)
        parser.parse()

    def insert_docs(self, docs: list[Document]) -> None: # look into also supporting embedding using ollama and other sources too
        '''
        Inserts a set of documents into the EigenDB instance. 
        The documents are vectorized using the provided model provider and model name.
        Throws an error if model_provivder = "none".
        Args:
            embeddings: A list of Embedding objects to be inserted.
        '''
        sentence_embeddings = self._vectorize_docs(docs)
        self.insert(embeddings=sentence_embeddings)

    def upsert(self, embeddings: list[Embedding]) -> None:
        '''
        Upserts a set of embeddings into the EigenDB instance.
        Args:
            embeddings: A list of Embedding objects to be upserted.
        '''
        res = put(
            url=self.url + '/embeddings/upsert',
            headers={
                'X-Eigen-API-Key': self.api_key
            },
            data=json.dumps({
                "embeddings": [e.to_dict() for e in embeddings]
            })
        )
        parser = ResponseParser(res)
        parser.parse()

    def upsert_docs(self, docs: list[Document]) -> None:
        '''
        Upserts a set of documents into the EigenDB instance.
        The documents are vectorized using the provided model provider and model name.
        Throws an error if model_provivder = "none".
        Args:
            embeddings: A list of Embedding objects to be upserted.
        '''
        sentence_embeddings = self._vectorize_docs(docs)
        self.upsert(embeddings=sentence_embeddings)

    def search(self, query: Embedding, k: int) -> dict[str, dict]:
        '''
        Performs a similarity search on the EigenDB instance using the provided query embedding and k value. 
        Args:
            query: An Embedding object representing the query vector.
            k: The number of nearest neighbors to retrieve.
        Returns:
            A dictionary mapping embedding IDs to their corresponding nearest neighbor information.
        '''
        res = get(
            url=self.url + '/embeddings/search',
            headers={
                'X-Eigen-API-Key': self.api_key
            },
            data=json.dumps({
                "queryVector": query.data,
                "k": k,
            })
        )
        parser = ResponseParser(res)
        parser.parse()

        return parser.data['nearest_neighbors']
    
    def search_docs(self, string: str, k: int) -> dict[str, dict]:
        '''
        Performs a similarity search on the EigenDB instance using the provided query string and k value.
        The string is vectorized using the provided model provider and model name.
        Throws an error if model_provivder = "none".
        Args:
            string: A string representing the query document.
            k: The number of nearest neighbors to retrieve.
        Returns:
            A dictionary mapping embedding IDs to their corresponding nearest neighbor information.
        '''
        sentence_embedding = self._vectorize_docs([Document(id=-1, data=string)])[0]
        return self.search(
            query=Embedding(
                id=-1,
                data=sentence_embedding.data
            ),
            k=k
        )

    def retrieve(self, ids: list[int]) -> dict[str, Embedding]:
        '''
        Retrieves embeddings from the EigenDB instance using the provided list of IDs.
        Args:
            ids: A list of embedding IDs to retrieve.
        Returns:
            A dictionary mapping embedding IDs to their corresponding Embedding objects.
        '''
        res = get(
            url=self.url + '/embeddings/retrieve',
            headers={
                'X-Eigen-API-Key': self.api_key
            },
            data=json.dumps({
                "ids": ids
            })
        )
        parser = ResponseParser(res)
        parser.parse()

        results: dict[str, Embedding] = {}
        for embedding in parser.data['embeddings']:
            results[str(embedding['id'])] = Embedding(
                id=embedding['id'],
                data=embedding['data'],
                metadata=embedding['metadata']
            )

        return results

    def delete(self, ids: list[int]) -> None:
        '''
        Deletes embeddings from the EigenDB instance using the provided list of IDs.
        Args:
            ids: A list of embedding IDs to delete.
        '''
        res = delete(
            url=self.url + '/embeddings/delete',
            headers={
                'X-Eigen-API-Key': self.api_key
            },
            data=json.dumps({
                "ids": ids
            })
        )
        parser = ResponseParser(res)
        parser.parse()

    def __repr__(self) -> str:
        return f"EigenDB_Index(url={self.url})"