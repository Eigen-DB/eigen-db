#!/usr/bin/python3

# WORK IN PROGRESS! 
# Working on gathering performance metrics 

import requests
import json
from time import time

BASE_URL = "http://127.0.0.1:8080"
API_KEY = "7e691a0cfae65c1254fb8f9879fa6e2d31e1700b5556500d39230a12a066fe22" # change this to proper API key

def setup() -> None:
    req_body = {
        "setOfComponents": [
            [3.2, -1.5],
            [4.7, 2.1],
            [-6.3, 3.4],
            [0.9, -4.8],
            [-2.7, 5.6],
            [1.3, -3.9],
            [2.4, 6.1],
            [-1.1, 3.0],
            [5.5, -2.2],
            [0.0, 4.4],
            [-3.6, -0.7],
            [4.1, 5.3],
            [-2.9, 2.8],
            [3.7, -3.6],
            [1.0, 0.5],
            [5.9, 1.7],
            [-4.4, -3.2],
            [2.8, 4.9],
            [-1.5, -2.4],
            [3.3, 1.6],
            [4.6, -1.3],
            [-2.1, 3.7],
            [1.8, -5.4],
            [3.9, 2.5],
            [-1.4, 4.2],
            [0.2, -3.1],
            [5.1, 1.3],
            [-2.8, -1.7],
            [3.0, 5.5],
            [1.5, -2.8],
            [-4.9, 3.1],
            [2.6, -4.5],
            [0.7, 3.8],
            [-3.3, 2.2],
            [4.0, -0.9],
            [-1.2, 4.9],
            [3.4, -2.6],
            [0.6, 1.8],
            [-2.5, -3.9],
            [5.3, 2.0],
            [-0.8, 3.3],
            [2.1, -4.2],
            [4.5, 1.4],
            [-3.7, -2.5],
            [1.9, 3.6],
            [0.3, -5.1],
            [4.8, -3.0],
            [-1.6, 2.9],
            [2.9, -4.0]
        ]
    }
    num_of_vectors = len(req_body['setOfComponents'])
    req = requests.put(
        url=BASE_URL + "/vector/bulk-insert",
        data=json.dumps(req_body),
        headers={
            "X-Eigen-API-Key": API_KEY
        }
    )

    if req.status_code != 200 or req.content.decode() != f"{num_of_vectors}/{num_of_vectors} vectors successfully inserted.":
        print(f"ERROR: request to /vector/bulk-insert failed. Response: {req.content.decode()}")
        exit(1)

'''
Returns avg times (seconds) to perform the indexing.
'''
def benchmark_indexing() -> float:
    times = []
    num_of_trials = 100

    req_body = {
        'queryVectorId': 1,
        'k': 5
    }

    for i in range(num_of_trials):
        start = time()
        req = requests.get(
            url=BASE_URL + "/vector/search",
            data=json.dumps(req_body),
            headers={
                "X-Eigen-API-Key": API_KEY
            }
        )
        end = time()

        if req.status_code != 200:
            print(f"ERROR: request to /vector/search failed. Response: {req.content.decode()}")
        else:
            times.append(end - start)
    
    return sum(times) / len(times)





if __name__ == "__main__":
    setup()

    mean = benchmark_indexing()
    print(mean)
    exit(0)