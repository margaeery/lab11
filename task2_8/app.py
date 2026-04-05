from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

app = FastAPI()


class DataResponse(BaseModel):
    data: dict


@app.get("/health")
def health():
    return {"status": "ok"}


@app.get("/")
def root():
    return {"message": "Hello, World!"}


@app.post("/data")
def data(body: dict):
    if not body:
        raise HTTPException(status_code=400, detail="Request body must not be empty")
    return body
