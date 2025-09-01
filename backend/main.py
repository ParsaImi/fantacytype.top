from fastapi import FastAPI
import os
import sys

current_dir = os.path.dirname(os.path.abspath(__file__))
# Add the parent directory to the Python path
sys.path.insert(0, current_dir)
from myroutes import sentences , auth
from dotenv import load_dotenv
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI()
load_dotenv()

# project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), '..'))
# sys.path.insert(0, project_root)

origins = [
    "http://localhost",
    "http://localhost:5173",
    "http://fantacytype.top:5173",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
class Settings:
    SECRET_KEY: str = os.getenv('SECRET_KEY' , "fallback_secret_key")

settings = Settings()


app.include_router(auth.router)
app.include_router(sentences.router)

@app.get("/")
async def root():
    return {"message": "Hellow nig"}
