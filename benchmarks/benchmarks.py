#!/usr/bin/python3

import requests
import json
import random
from time import time

EIGEN_ENDPOINT = "http://127.0.0.1:8080"
API_KEY = "20417e0c21cea44e4b4fc90a06f57658" # change this to proper API key

####### HELPER FUNCTIONS #######

def create_random_vector(dim: int) -> "list[float]":
    v = []
    for i in range(dim):
        v.append(random.uniform(0.0, 200.0))
    return v

def setup(num_vectors: int, dim: int) -> None:
    req_body = {"vectors": []}
    for i in range(num_vectors):
        req_body["vectors"].append({
            "embedding": create_random_vector(dim),
            "id": i+1
        })

    req = requests.put(
        url=EIGEN_ENDPOINT + "/vector/bulk-insert",
        data=json.dumps(req_body),
        headers={
            "X-Eigen-API-Key": API_KEY
        }
    )

    if req.status_code != 200 or json.loads(req.content.decode()).get("message") != f"{num_vectors}/{num_vectors} vectors successfully inserted.":
        print(f"ERROR: request to /vector/bulk-insert failed. Response: {req.content.decode()}")
        exit(1)


####### BENCHMARKING FUNCTIONS ####### 

'''
Returns mean time in seconds to perform similarity search.
'''
def benchmark_indexing(num_vectors: int, num_trials: int, k: int) -> float:
    times_secs = []

    req_body = {
        'queryVectorId': random.randint(1, num_vectors),
        'k': k
    }

    for i in range(num_trials):
        start = time()
        req = requests.get(
            url=EIGEN_ENDPOINT + "/vector/search",
            data=json.dumps(req_body),
            headers={
                "X-Eigen-API-Key": API_KEY
            }
        )
        end = time()

        if req.status_code != 200:
            print(f"ERROR: request to /vector/search failed. Response: {req.content.decode()}")
        else:
            times_secs.append(end - start)
    
    return sum(times_secs) / len(times_secs)

'''
Returns the mean time in seconds to insert a vector
'''
def benchmark_inserting(num_vectors: int, dim: int, num_trials: int) -> float:
    times_secs = []
    req_body = {
        "vector": {
            "embedding": create_random_vector(dim),
            "id": None # set dynamically during trials
        }
    }

    for i in range(num_trials):
        start = time()
        req_body["vector"]["id"] = num_vectors + i + 1
        req = requests.put(
            url=EIGEN_ENDPOINT + "/vector/insert",
            data=json.dumps(req_body),
            headers={
                "X-Eigen-API-Key": API_KEY
            }
        )
        end = time()

        if req.status_code != 200:
            print(f"ERROR: request to /vector/search failed. Response: {req.content.decode()}")
        else:
            times_secs.append(end - start)

    return sum(times_secs) / len(times_secs)




if __name__ == "__main__":
    #setup(
    #    num_vectors=10_000, 
    #    dim=512
    #)

    index_mean = benchmark_indexing(
        num_vectors=10_000,
        num_trials=100, 
        k=50
    )
    print(f"index mean: {index_mean}s")

    insert_mean = benchmark_inserting(
        num_vectors=10_000,
        num_trials=100,
        dim=512
    )
    print(f"insert mean: {insert_mean}s")

    exit(0)