import os
import json
from fastapi import APIRouter
from fastapi import Header
from fastapi import HTTPException
from huggingface_hub import InferenceClient
from typing import List

from helpers import vectors
from helpers.model import perform_inference
from schemas.data import Data
from schemas import eigen_db_res

router = APIRouter()
inference_client = InferenceClient(
    model=os.getenv('HUGGING_FACE_EMBEDDING_MODEL'),
    token=os.getenv('HUGGING_FACE_API_TOKEN'),
)

@router.put("/bulk-upload")
async def bulk_upload_data(datas: List[Data], x_eigen_api_key: str = Header(None)) -> eigen_db_res.Response:
    model_type = os.getenv('MODEL_TYPE')
    embeddings = []
    try:
        for d in datas:
            embedding = perform_inference(
                inference_client=inference_client,
                model_type=model_type,
                data=d.data
            )
            embeddings.append(embedding)
    except Exception as e:
        raise HTTPException(
            status_code=400,
            detail={"error": str(e)}
        )

    res = vectors.bulk_insert(
        vectors=embeddings,
        ids=[d.id for d in datas],
        api_key=x_eigen_api_key
    )

    if res.status_code != 200:
        raise HTTPException(
            status_code=res.status_code,
            detail=json.loads(res.content.decode())
        )
    return json.loads(res.content.decode())