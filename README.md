# Eigen DB

### The blazingly fast in-memory vector database

Payload to `/vector/bulk-insert` to create test data (must use 2D vector space) and a [visual](https://www.desmos.com/calculator/pjjei9akcx):

Keep in mind that in a real application, vector dimensionality is usually much larger depending on the complextiy and granularity of the data.
```json
{
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
```