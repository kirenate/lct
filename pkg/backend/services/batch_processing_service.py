from dataclasses import dataclass
import uuid
import io

from sqlalchemy.dialects.postgresql.pg_catalog import format_type

from repositories.pg_repository import PgRepository
from repositories.redis_repository import RedisRepository
from repositories.minio_repository import MinioRepository
from pdf2image import convert_from_bytes
import pymupdf
from PIL  import Image

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
        ext = filename.split(".")[-1]
        if ext not in [ "png", "jpg", "jpeg", "tiff"]:
            raise Exception("Invalid file extension")
        # if ext == "pdf":
        #     img = io.BytesIO()
        #     convert_from_bytes(batch)[0].save(fp=img, format="png")
        #     self.minio.client.upload_fileobj(img, self.minio.bucket_name, str(uid) + ".png")
        #     return uid

        self.minio.client.upload_fileobj(io.BytesIO(batch), self.minio.bucket_name, str(uid)+".png")
        return uid

    async def delete_document(self, document_id:uuid.UUID)-> None:
        self.minio.client.delete_object(Bucket=self.minio.bucket_name, Key=str(document_id))

