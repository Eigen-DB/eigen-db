from eigen_client.client import Client
from eigen_client.data_types import Document
from eigen_client.index import Index
import pandas as pd

def create_index_if_not_exist(client: Client) -> None:
    try:
        client.get_index_stats('tweets')
    except Exception: # index does not exist
        print('Index does not exist... creating it now.')
        index = client.create_index_from_model(
            index_name='tweets',
            model_provider='ollama',
            model_name='all-minilm:22m'
        )
        print('Upserting tweets...')
        df = pd.read_csv('data/tweets.csv')
        tweets = [Document(id=index, data=row['text'], metadata={"id": str(row['ids']), "content": row['text'], "user": row['user']}) for index, row in df.iterrows()]
        index.upsert_docs(tweets)
        print('Upsert complete.')

def keyword_search(search_query: str, desired_results: int) -> list[dict]:
    df = pd.read_csv('data/tweets.csv')
    words = search_query.lower().strip().split(' ')
    pattern = r'\b(' + '|'.join(words) + r')\b'
    results = df[df['text'].str.contains(pattern, case=False, regex=True)].head(desired_results)
    output = []
    i = 0
    for index, row in results.iterrows():
        output.append({
            "metadata": {
                "id": str(row['ids']),
                "content": row['text'],
                "user": row['user']
            },
            "rank": i
        })
        i += 1
    return output

def semantic_search(index: Index, search_query: str, desired_results: int) -> list[dict]:
    results = index.search_docs(string=search_query, k=desired_results)
    recommendations = [results[reco_id] for reco_id in results.keys()]
    return recommendations
