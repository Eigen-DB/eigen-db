# EigenDB's Official Python API 🐍

A Python wrapper for EigenDB's REST API. 

### Example usage:
```py
import os
from eigen_py.index import Index
from eigen_py.data_types import Document

index = Index(
    url="http://localhost:8080",
    api_key="your eigendb api key...",
    model_name="text-embedding-3-small",
    model_provider="openai",
    model_provider_api_key="your openai api key..."
)

documents = [
    Document(id=1, data="Fresh herbs boost flavor.", metadata={"recipe_id": "123"}),
    Document(id=2, data="Slow simmer blends soup.", metadata={"recipe_id": "456"}),
    Document(id=3, data="Homemade bread smells great.", metadata={"recipe_id": "789"}),
    Document(id=4, data="Grilled veggies taste sweeter.", metadata={"recipe_id": "987"}),
    Document(id=5, data="Cast iron sears steak well.", metadata={"recipe_id": "654"})
]

index.upsert_docs(documents)

results = index.search_docs(
    string="Baking",
    k=3
)

print(results)
```

Made with ❤️ by developers, for developers.