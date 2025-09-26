from flask import Flask, request
from flask_cors import CORS
from dotenv import load_dotenv
from eigen_client.index import Index
from eigen_client.client import Client
import os

from utils import create_index_if_not_exist, keyword_search, semantic_search

load_dotenv()

app = Flask(__name__)
CORS(app)
index = Index(
    url=os.getenv('EIGENDB_INSTANCE_URL'),
    api_key=os.getenv('EIGENDB_API_KEY'),
    index_name='tweets',
    model_provider='ollama',
    model_name='all-minilm:22m'
)

@app.route('/health', methods=['GET'])
def health():
    return {'status': 'ok'}, 200

@app.route('/recommend', methods=['POST']) # returns list of metadatas and similarity ranks
def recommend():
    body = request.get_json()
    use_eigendb = body.get('use_eigendb', True)
    tweet = body.get('tweet')
    if tweet:
        desired_results = body.get('desired_results', 5)
        if not use_eigendb:
            recommendations = keyword_search(tweet, desired_results)
            return {'recommendations': recommendations}, 200
        
        recommendations = semantic_search(index, tweet, desired_results)
        return {'recommendations': recommendations}, 200
    return {'error': 'Invalid request body'}, 400

if __name__ == '__main__':
    client = Client(
        url=os.getenv('EIGENDB_INSTANCE_URL'),
        api_key=os.getenv('EIGENDB_API_KEY')
    )
    create_index_if_not_exist(client)
    app.run(debug=False) # change in prod 