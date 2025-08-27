from requests import get, put, delete
from typing import Literal
import json
import ollama
import openai
import tiktoken

from eigen_py.response import ResponseParser
from eigen_py.data_types import Embedding, Text
from eigen_py.supported_models import SUPPORTED_MODEL_PROVIDERS, SUPPORTED_MODELS

#logging.set_verbosity_error()

class EigenIndex:
    def __init__(self, 
        url: str,
        api_key: str,
        model_provider: Literal["openai", "ollama_local", "ollama_cloud", "none"] = "openai",
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
        
        self.test_auth()
        if self.model_provider == "openai":
            self.openai_client = openai.Client(api_key=self.model_provider_api_key)
        elif self.model_provider == "ollama_local":
            self.ollama_client = ollama
        elif self.model_provider == "ollama_cloud":
            self.ollama_client = ollama.Client(
                host="https://ollama.com",
                headers={
                    "Authorization": self.model_provider_api_key
                }
            )

    def test_auth(self) -> bool:
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
    
    def vectorize_text(self, texts: list[Text]) -> list[Embedding]:
        output_embeddings: list[Embedding] = []
        if self.model_provider in ["ollama_local", "ollama_cloud"]: # TEST OLLAMA TURBO
            response = self.ollama_client.embed(
                model=self.model['name'],
                input=[text.data for text in texts],
                truncate=True
            )
            for i in range(len(response.embeddings)):
                output_embeddings.append(
                    Embedding(
                        id=texts[i].id,
                        data=response.embeddings[i],
                        metadata=texts[i].metadata
                    )
                )
        elif self.model_provider == "openai":
            enc = tiktoken.get_encoding("cl100k_base") # should work for all openai embedding models we support at the moment: https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
            token_limit = self.model['metadata']['token_limit']
            for text in texts:
                encoded_text = enc.encode(text.data)
                num_tokens = len(encoded_text)
                if num_tokens > token_limit:
                    print(f"Warning: Text with ID {text.id} exceeds the token limit of {token_limit} tokens for model {self.model['name']}. It has {num_tokens} tokens and will be truncated.")
                    truncated_encoded_text = encoded_text[:token_limit]
                    text.data = enc.decode(truncated_encoded_text)

            response = self.openai_client.embeddings.create(
                model=self.model['name'],
                input=[text.data for text in texts],
            )
            for i in range(len(response.data)):
                output_embeddings.append(
                    Embedding(
                        id=texts[i].id,
                        data=response.data[i].embedding,
                        metadata=texts[i].metadata
                    )
                )
        elif self.model_provider == "none":
            raise Exception("No model provider specified. Please set a model provider to vectorize text.")
        else:
            raise Exception(f"Invalid model provider: {self.model_provider}. Supported providers are: {', '.join(SUPPORTED_MODEL_PROVIDERS)}.")

        return output_embeddings

    def insert(self, embeddings: list[Embedding]) -> None:
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

    def insert_text(self, texts: list[Text]) -> None: # look into also supporting embedding using ollama and other sources too
        sentence_embeddings = self.vectorize_text(texts)
        self.insert(embeddings=sentence_embeddings)

    def upsert(self, embeddings: list[Embedding]) -> None:
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

    def upsert_text(self, texts: list[Text]) -> None:
        sentence_embeddings = self.vectorize_text(texts)
        self.upsert(embeddings=sentence_embeddings)

    def search(self, query: Embedding, k: int) -> dict[str, dict]:
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
    
    def search_text(self, text: str, k: int) -> dict[str, dict]:
        sentence_embedding = self.vectorize_text([Text(id=-1, data=text)])[0]
        return self.search(
            query=Embedding(
                id=-1,
                data=sentence_embedding.data
            ),
            k=k
        )

    def retrieve(self, ids: list[int]) -> dict[str, Embedding]:
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
        return f"EigenIndex(url={self.url})"