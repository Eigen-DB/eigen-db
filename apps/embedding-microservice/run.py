import requests
import os
import json
from fastapi import FastAPI
from fastapi import APIRouter
from fastapi import Request
from fastapi.responses import JSONResponse
from starlette.responses import StreamingResponse as Response
from routes import upload
from routes import bulk_upload
from typing import Callable
from typing import Awaitable
from dotenv import load_dotenv

load_dotenv()

app = FastAPI()

# defining the middleware
@app.middleware("http")
async def validate_api_key(request: Request, call_next: Callable[[Request], Awaitable[Response]]) -> Response: # look into header schema
    eigen_api_key = request.headers.get('X-Eigen-API-Key')
    eigen_res = requests.get(
        url=os.getenv('EIGEN_DB_INSTANCE_URL') + '/test-auth',
        headers={"X-Eigen-API-Key": eigen_api_key}
    )
    if eigen_res.status_code != 200: # invalid API key
        return JSONResponse(
            content=json.loads(eigen_res.content.decode()),
            status_code=eigen_res.status_code
        )

    response = await call_next(request)
    return response

#@app.middleware("http")
#async def format_response(request: Request, call_next: Callable[[Request], Awaitable[Response]]) -> Response:
#    response = await call_next(request) # forward request to endpoint
#
#    res_body = b"".join([chunk async for chunk in response.body_iterator])
#    res_body: dict = json.loads(res_body)
#
#    if res_body.get('details') != None:
#        res_body = res_body.get('details')
#
#    if 'error' in res_body.keys() and res_body.get('error') == None: # get rid of {... error: null, ...} from response
#        res_body.pop('error')
#
#    response.body = res_body
#    return response

api_v1_router = APIRouter(prefix="/api/v1")
api_v1_router.include_router(upload.router)
api_v1_router.include_router(bulk_upload.router)

app.include_router(api_v1_router)