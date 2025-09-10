#!/usr/bin/python3

import yaml
import os

CONFIG_LOCATION = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'eigen', 'config.yml')
print(CONFIG_LOCATION)
SUPPORTED_MODELS = [
    {
        "name": "all-minilm:22m",
        "supported_provider": "ollama",
        "dimensions": 384,
        "metric": "cosine",
        "metadata": {}
    },
    {
        "name": "nomic-embed-text:v1.5",
        "supported_provider": "ollama",
        "dimensions": 768,
        "metric": "cosine",
        "metadata": {}
    },
    {
        "name": "mxbai-embed-large:335m",
        "supported_provider": "ollama",
        "dimensions": 1024,
        "metric": "cosine",
        "metadata": {}
    },
    {
        "name": "snowflake-arctic-embed2:568m",
        "supported_provider": "ollama",
        "dimensions": 1024,
        "metric": "cosine",
        "metadata": {}
    },
    {
        "name": "snowflake-arctic-embed:335m",
        "supported_provider": "ollama",
        "dimensions": 1024,
        "metric": "cosine",
        "metadata": {}
    },
    {
        "name": "bge-m3:567m",
        "supported_provider": "ollama",
        "dimensions": 1024,
        "metric": "cosine",
        "metadata": {}
    },
    {
        "name": "bge-large:335m",
        "supported_provider": "ollama",
        "dimensions": 1024,
        "metric": "cosine",
        "metadata": {}
    },
    {
        "name": "text-embedding-3-small",
        "supported_provider": "openai",
        "dimensions": 1536,
        "metric": "cosine",
        "metadata": {
            "token_limit": 8192 
        }
    },
    {
        "name": "text-embedding-3-large",
        "supported_provider": "openai",
        "dimensions": 3072,
        "metric": "cosine",
        "metadata": {
            "token_limit": 8192 
        }
    },
    {
        "name": "text-embedding-ada-002",
        "supported_provider": "openai",
        "dimensions": 1536,
        "metric": "cosine",
        "metadata": {
            "token_limit": 8192 
        }
    },
]

def generate_embedding_menu() -> None:
    print("Select an embedding model:")
    for i in range(len(SUPPORTED_MODELS)):
        print(f"({i+1}) [{SUPPORTED_MODELS[i]['supported_provider']}] {SUPPORTED_MODELS[i]['name']}")
    print(f"({len(SUPPORTED_MODELS)+1}) None of the above / Custom embedding model")

def generate_config(selected_model_index: int, dimensions: int = None, metric: str = None) -> dict:
    if not (dimensions and metric):
        selected_model = SUPPORTED_MODELS[selected_model_index]
    return {
        'persistence': {
            'timeInterval': '3s',
        },
        'api': {
            'port': 8080,
            'address': '0.0.0.0'
        },
        'indexConfig': {
            'dimensions': dimensions if dimensions else selected_model['dimensions'],
            'similarityMetric': metric if metric else selected_model['metric'],
        }
    }

def main() -> None:
    generate_embedding_menu()
    while True:
        try:
            selected_model_index = int(input(f"Selected model (1-{len(SUPPORTED_MODELS)+1}): "))-1
            if selected_model_index > len(SUPPORTED_MODELS) or selected_model_index < 0:
                raise ValueError
            
            if selected_model_index == len(SUPPORTED_MODELS):
                while True:
                    try:
                        dimensions = int(input("Enter the number of dimensions for your custom embedding model: "))
                        if dimensions <= 1:
                            raise ValueError
                        break
                    except ValueError:
                        print("Invalid input. Please enter a valid integer for dimensions (>1).")
                        continue

                while True:    
                    metric = input("Enter the distance metric ('cosine', 'euclidean', 'ip'): ").strip().lower()
                    if metric not in ['cosine', 'l2', 'ip']:
                        print("Invalid metric. Please enter 'cosine', 'euclidean', or 'ip'.")
                        continue
                    break
                cfg = generate_config(selected_model_index, dimensions, metric)
            else:
                cfg = generate_config(selected_model_index)

            with open(CONFIG_LOCATION, 'w') as file:
                yaml.safe_dump(cfg, file)

            print(f'Configuration file successfully created for {SUPPORTED_MODELS[selected_model_index]["name"]}.')   
            break             
        except ValueError:
            print(f"Invalid input. Please enter one of 1-{len(SUPPORTED_MODELS)+1}.")
            continue


if __name__ == '__main__':
    main()