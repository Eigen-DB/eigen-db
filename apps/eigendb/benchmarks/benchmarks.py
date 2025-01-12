#!/usr/bin/python3

import requests
import json
import random
import os
import subprocess
import psutil
import pandas as pd
from time import time, sleep

EIGEN_ENDPOINT = "http://127.0.0.1:8080"
API_KEY = "20417e0c21cea44e4b4fc90a06f57658" # change this to proper API key
TEST_PARAMS = [
    {
        "num_vectors": 500,
        "num_trials": 100,
        "k": 10,
        "dim": 512,
    },
    {
        "num_vectors": 500,
        "num_trials": 100,
        "k": 50,
        "dim": 512,
    },
    {
        "num_vectors": 500,
        "num_trials": 100,
        "k": 100,
        "dim": 512,
    },
    {
        "num_vectors": 500,
        "num_trials": 100,
        "k": 500,
        "dim": 512,
    },
    {
        "num_vectors": 1_000,
        "num_trials": 100,
        "k": 10,
        "dim": 512,
    },
    {
        "num_vectors": 1_000,
        "num_trials": 100,
        "k": 50,
        "dim": 512,
    },
    {
        "num_vectors": 1_000,
        "num_trials": 100,
        "k": 100,
        "dim": 512,
    },
    {
        "num_vectors": 1_000,
        "num_trials": 100,
        "k": 500,
        "dim": 512,
    },
    {
        "num_vectors": 5_000,
        "num_trials": 100,
        "k": 10,
        "dim": 512,
    },
    {
        "num_vectors": 5_000,
        "num_trials": 100,
        "k": 50,
        "dim": 512,
    },
    {
        "num_vectors": 5_000,
        "num_trials": 100,
        "k": 100,
        "dim": 512,
    },
    {
        "num_vectors": 5_000,
        "num_trials": 100,
        "k": 500,
        "dim": 512,
    },
    {
        "num_vectors": 10_000,
        "num_trials": 100,
        "k": 10,
        "dim": 512,
    },
    {
        "num_vectors": 10_000,
        "num_trials": 100,
        "k": 50,
        "dim": 512,
    },
    {
        "num_vectors": 10_000,
        "num_trials": 100,
        "k": 100,
        "dim": 512,
    },
    {
        "num_vectors": 10_000,
        "num_trials": 100,
        "k": 500,
        "dim": 512,
    },
]

####### HELPER FUNCTIONS #######

def create_random_vector(dim: int) -> "list[float]":
    v = []
    for i in range(dim):
        v.append(random.uniform(-200.0, 200.0))
    return v

# reuse_req_body allows multiple setups with the same num_vectors to use the first setup's randomly generated vectors instead of re-computing them.
def setup(num_vectors: int, dim: int, reuse_req_body: dict = None) -> None:
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
        print(f"SETUP ERROR: request to /vector/bulk-insert failed. Response: {req.content.decode()}")

def start_eigen_db() -> None:
    log_stream = open(f'./eigen_db-{time()}.log', mode='a')
    subprocess.Popen(
        args=['./eigen_db'],
        cwd='../',
        stdout=log_stream,
        stderr=log_stream
    )

def find_and_kill_process(binary_name: str) -> None:
    for process in psutil.process_iter(['pid', 'name']):
        try:
            if binary_name == process.info['name']:
                print(f"Killing process {binary_name} with PID {process.info['pid']}")
                process.kill()
                return
        except (psutil.NoSuchProcess, psutil.AccessDenied, psutil.ZombieProcess):
            pass
    print(f"No running process found with name \"{binary_name}\"")

def cleanup(restart=True) -> None:
    find_and_kill_process('eigen_db')
    sleep(1) # ensure the process is truly killed

    if os.path.exists("../eigen/hnsw_index.bin"):
        os.remove("../eigen/hnsw_index.bin")  

    if os.path.exists("../eigen/vector_space.vec"):
        os.remove("../eigen/vector_space.vec")    

    if restart:
        print('Restarting EigenDB...')
        start_eigen_db()


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
    print('WARNING: running this script will delete all existing persisted data. This script it mean to be ran for development purposes.')
    answer = input('Are you sure you want proceed? (y/N): ')
    if answer != 'y':
        print('Exitting...')
        exit(0)

    index_df = pd.DataFrame(
        columns=['num_vectors', 'k', 'dim', 'mean_time_secs']
    )

    insert_df = pd.DataFrame(
        columns=['num_vectors', 'k', 'dim', 'mean_time_secs']
    )

    cleanup(restart=False)
    start_eigen_db()

    for param in TEST_PARAMS:
        print(f"Metrics for the following params:\n{json.dumps(param, indent=4)}")

        setup(
            num_vectors=param["num_vectors"], 
            dim=param["dim"]
        )

        index_mean = benchmark_indexing(
            num_vectors=param["num_vectors"],
            num_trials=param["num_trials"], 
            k=param["k"]
        )
        index_df.loc[len(index_df)] = [param['num_vectors'], param['k'], param['dim'], index_mean] # append row in the df
        print(f"index mean: {index_mean}s")

        insert_mean = benchmark_inserting(
            num_vectors=param["num_vectors"],
            num_trials=param["num_trials"],
            dim=param["dim"]
        )
        insert_df.loc[len(insert_df)] = [param['num_vectors'], param['k'], param['dim'], insert_mean] # append row in the df
        print(f"insert mean: {insert_mean}s")

        cleanup()
    
    index_df.to_csv(f'./indexing_mean_{time()}.csv', index=False)
    insert_df.to_csv(f'./inserting_mean_{time()}.csv', index=False)

    cleanup(restart=False)
    exit(0)