import multiprocessing as mp

from pydantic_settings import BaseSettings, SettingsConfigDict
from schemas.base import PureBaseModel


class Postgres(PureBaseModel):
    database: str = "db_main"
    host: str = "localhost"
    port: int = 5432
    username: str = "db_main"
    password: str = "db_main"


class Uvicorn(PureBaseModel):
    host: str = "localhost"
    port: int = 8000
    workers: int = mp.cpu_count() * 2 + 1
    log_level: str = "WARNING"


class Redis(PureBaseModel):
    host: str = "localhost"
    port: int = 6379


class AppSettings(BaseSettings):
    pg: Postgres = Postgres()
    uvicorn: Uvicorn = Uvicorn()
    redis: Redis = Redis()

    model_config = SettingsConfigDict(env_file=".env", env_prefix="_", env_nested_delimiter="__")


app_settings = AppSettings()
