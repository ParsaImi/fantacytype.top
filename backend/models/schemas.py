from datetime import datetime, timezone
from fastapi import FastAPI, HTTPException, Depends, status
from pydantic import BaseModel
from sqlalchemy.orm import declarative_base
from sqlalchemy.sql import func
from sqlalchemy import DateTime, create_engine, Column, Integer, String
from sqlalchemy.orm import sessionmaker , Session
from database.database import get_engine
engine = get_engine()
# models
Base = declarative_base()

class User(Base):
    __tablename__ = "users"
    id = Column(Integer, primary_key=True, index=True)
    email = Column(String, unique=True, index=True)
    username = Column(String, unique=True, index=True)
    hashed_password = Column(String)
    created_on = Column(DateTime, default= datetime.now(timezone.utc), nullable=False)


Base.metadata.create_all(bind=engine)
# Pydantic models
class UserLogin(BaseModel):
    username: str
    password: str

class UserCreate(BaseModel):
    email: str
    username: str
    password: str

class UserResponse(BaseModel):
    id: int
    email: str
    username: str
    created_on: datetime
    
#  Helper Function

