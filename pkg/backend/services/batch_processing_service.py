from dataclasses import dataclass
import uuid
import io
from repositories.pg_repository import PgRepository
from repositories.redis_repository import RedisRepository
from repositories.minio_repository import MinioRepository

@dataclass
class Service:
    redis: RedisRepository
    pg: PgRepository
    minio: MinioRepository

    def __init__(self):
        self.redis = RedisRepository()
        self.pg = PgRepository()
        self.minio = MinioRepository()

    async def save_document(self, batch: bytes, filename: str) ->  uuid.UUID :
        uid = uuid.uuid4()
        self.minio.client.upload_fileobj(io.BytesIO(batch), self.minio.bucket_name, str(uid))
        return uid

    async def delete_document(self, document_id:uuid.UUID)-> None:
        self.minio.client.delete_object(Bucket=self.minio.bucket_name, Key=str(document_id))

