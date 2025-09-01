import os
from typing import Annotated, Optional
from fastapi import security
import jwt
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from jwt.exceptions import ExpiredSignatureError  # Add this import
from datetime import timedelta, datetime
from fastapi import APIRouter, HTTPException, status
from sqlalchemy.orm import Session
from starlette.status import HTTP_201_CREATED
from database.database import SessionLocal
from models.schemas import UserCreate, User, UserResponse , UserLogin
from passlib.context import CryptContext
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from pydantic import BaseModel, Field
from dotenv import load_dotenv
from fastapi import Depends

load_dotenv()


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

router = APIRouter(
        prefix="/auth",
        tags=["auth"]
)

# Password Hashing 
pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="auth/login")
ACCESS_TOKEN_EXPIRE_MINUTES = 30
SECRET_KEY = "c4688488d8b6a51e00426bb9e8f15d0c7566e485c932f6b999f5a067e5d4e9a8"
ALGORITHM = "HS256"



# Token model
class Token(BaseModel):
    access_token: str
    token_type: str

class PasswordReset(BaseModel):
    current_password: str
    new_password: str = Field(... , min_length=8)

class TokenData(BaseModel):
    username : str | None  = None


def create_access_token(data: dict , expires_delta):
    to_encode = data.copy()
    if expires_delta:
        expire = datetime.now() + expires_delta
    else:
        expire = datetime.now() + timedelta(minutes=15)

    to_encode.update({"exp" : expire})
    encoded_jwt = jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)
    return encoded_jwt

def get_user_by_email(db : Session, email: str):
    return db.query(User).filter(User.email == email).first()


def get_user_by_username(db : Session, username: str):
    return db.query(User).filter(User.username == username).first()


def get_current_user(token : Annotated[str, Depends(oauth2_scheme)], db : Annotated[Session , Depends(get_db)]):
    try:
        payload = jwt.decode(token, "hereisthekey", algorithms=[ALGORITHM])
        username = payload.get("sub")
        if username is None :
            raise HTTPException(status_code=404 , detail= "User not found!")
        token_data = TokenData(username=username)
    except:
        raise HTTPException(status_code=404 , detail= "User not found!")
    
    user = get_user_by_username(db , username=token_data.username)
    if user is None:
        raise HTTPException(status_code=404 , detail= "User not found!")
    return user


def get_hashed_password(password: str):
    return pwd_context.hash(password)



def verify_password(plain_password, hashed_password):
    return pwd_context.verify(plain_password, hashed_password)


@router.post("/signup" , response_model=UserResponse, status_code=HTTP_201_CREATED)
def signup(user: UserCreate , db : Annotated[Session , Depends(get_db)]):
    if get_user_by_email(db , user.email):
        raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Email already in use"
        )
    if get_user_by_username(db , user.username):
        raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Username already in use"
    )
    db_user = User(
            email = user.email,
            username = user.username,
            hashed_password=get_hashed_password(user.password)
    )

    db.add(db_user)
    db.commit()
    db.refresh(db_user)

    return db_user


@router.post("/verify")
def verify(credentials : HTTPAuthorizationCredentials):
    try:
        token = credentials.credentials
        payload = jwt.decode(token, SECRET_KEY, algorithms=ALGORITHM)
        return {"verify" : True , "username": payload.get("sub"), "id" : payload.get("id")}
    except ExpiredSignatureError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Token has expired",
            headers={"WWW-Authenticate": "Bearer"},
        )
    except jwt.PyJWTError as e:
        print(e)
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Could not validate credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )

@router.post("/login", response_model=Token)
def login(form_data : Annotated[OAuth2PasswordRequestForm, Depends()], db : Annotated[Session , Depends(get_db)]):
    print(form_data.username + " ", form_data.password)
    user = get_user_by_email(db , form_data.username)
    if not user :
        user = get_user_by_username(db , form_data.username)
    if not user or not verify_password(form_data.password, user.hashed_password):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect username/email or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    print(str(user.username))
    print(type(user.username))
    access_token = create_access_token(
            data={"id": user.id, "sub": str(user.username)},
            expires_delta=access_token_expires
    ) 
    return {"access_token": access_token, "token_type": "bearer"}



@router.put("/update-password", status_code=HTTP_201_CREATED)
def update_password(
        password_data: PasswordReset,
        current_user: Annotated[User , Depends(get_current_user)],
        db: Annotated[Session , Depends(get_db)]
):
    if not verify_password(password_data.current_password , current_user.hashed_password):
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Current password is incorrect"
        )
    current_user.hashed_password = get_hashed_password(password_data.new_password)
    db.add(current_user)
    db.commit()
    db.refresh(current_user)

    



