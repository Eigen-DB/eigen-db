#!/usr/bin/python3

import http.client
import json

def health_check() -> None: 
    conn = http.client.HTTPConnection("127.0.0.1:8080")

    conn.request("GET", "/health")
    response = conn.getresponse()
    data = json.loads(response.read().decode())
    conn.close()

    if response.status == 200 and data["status"] == "healthy":
        exit(0)
    else:
        exit(1)

if __name__ == "__main__":
    health_check()