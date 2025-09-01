from contextlib import contextmanager
import os
import sqlite3
from typing import Optional, Generator
from fastapi import APIRouter, HTTPException, status , Query
from enum import Enum
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from pathlib import Path
router = APIRouter(
        prefix="/corpus",
        tags=["corpus"]
)

BASE_PATH = Path(__file__).resolve().parent.parent

# Path to your SQLite database file
DB_PATHS = { "en_US" : BASE_PATH / "data" / "real_typing_sentences.db" , 
            "fa_IR" : BASE_PATH / "data" / "real_typing_sentences_fa.db"
}

class SupportedLocales(str , Enum):
    english = "en_US"
    persian = "fa_IR"

class SentenceResponse(BaseModel):
    sentence : str
    locale : str

@contextmanager
def get_db_connection(db_path: str):
    if not os.path.exists(db_path):
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Database file not found: {db_path}"
        )
    conn = None

    try:
        conn = sqlite3.connect(db_path)
        yield conn
    except sqlite3.Error as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Database error: {str(e)}"
        )
    finally:
        if conn:
            conn.close()


def get_random_sentence(locale: str) -> Optional[str]:
    db_path = DB_PATHS.get(locale)
    if not db_path:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=f"Unsupported locale: {locale}"
        )
    try:
        with get_db_connection(db_path) as conn:
            cursor = conn.cursor()
            cursor.execute("SELECT content FROM typing_sentences ORDER BY RANDOM() LIMIT 1;")
            result = cursor.fetchone()
            return result[0] if result else None
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to fetch sentence: {str(e)}"
        )




@router.get("/sentence", response_model=SentenceResponse)
def generate_single_sentence(
    locale: SupportedLocales = Query(SupportedLocales.english, description="Language locale")):
    sentence = get_random_sentence(locale.value)
    if not sentence:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"No sentences found for locale: {locale.value}"
        )

    return SentenceResponse(
        sentence=sentence,
        locale=locale.value
    )
   
