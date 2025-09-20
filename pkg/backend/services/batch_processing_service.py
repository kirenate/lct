from dataclasses import dataclass
import repositories.minio_repository
import repositories.pg_repository
import repositories.redis_repository
from botocore.exceptions import ClientError
import io
from loguru import logger
from repositories.pg_repository import PgRepository
from repositories.redis_repository import RedisRepository
from repositories.minio_repository import MinioRepository
import os

@dataclass
class Service:
    redis: RedisRepository
    pg: PgRepository
    minio: MinioRepository

    def __init__(self):
        self.redis = RedisRepository()
        self.pg = PgRepository()
        self.minio = MinioRepository()

    async def save_batch(self, batch: bytes, filename: str) -> ClientError | None :
        object_name = os.path.basename(filename)
        try:
            self.minio.client.upload_fileobj(io.BytesIO(batch), self.minio.bucket_name, object_name)

        except ClientError as e:
            return e
        return
