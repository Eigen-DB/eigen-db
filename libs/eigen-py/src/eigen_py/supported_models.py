SUPPORTED_MODEL_PROVIDERS = ["openai", "ollama_local", "ollama_cloud", "none"]
SUPPORTED_MODELS = [
    {
        "name": "snowflake-arctic-embed2:568m",
        "supported_providers": ["ollama_local", "ollama_cloud"],
        "dimensions": 1024,
        "metric": "cosine",
        "metadata": {}
    },
    {
        "name": "snowflake-arctic-embed:335m",
        "supported_providers": ["ollama_local", "ollama_cloud"],
        "dimensions": 1024,
        "metric": "cosine",
        "metadata": {}
    },
    {
        "name": "bge-m3:567m",
        "supported_providers": ["ollama_local", "ollama_cloud"],
        "dimensions": 1024,
        "metric": "cosine",
        "metadata": {}
    },
    {
        "name": "bge-large:335m",
        "supported_providers": ["ollama_local", "ollama_cloud"],
        "dimensions": 1024,
        "metric": "cosine",
        "metadata": {}
    },
    {
        "name": "text-embedding-3-small",
        "supported_providers": ["openai"],
        "dimensions": 1536,
        "metric": "cosine",
        "metadata": {
            "token_limit": 8192 
        }
    },
    {
        "name": "text-embedding-3-large",
        "supported_providers": ["openai"],
        "dimensions": 3072,
        "metric": "cosine",
        "metadata": {
            "token_limit": 8192 
        }
    },
    {
        "name": "text-embedding-ada-002",
        "supported_providers": ["openai"],
        "dimensions": 1536,
        "metric": "cosine",
        "metadata": {
            "token_limit": 8192 
        }
    },
]