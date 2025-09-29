from dataclasses import dataclass
import uuid
import io

from sqlalchemy.dialects.postgresql.pg_catalog import format_type

import schemas.responses
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

    async def save_document(self, batch: bytes, page_name: str, document_id:uuid.UUID) ->  uuid.UUID :
        uid = uuid.uuid4()
        spl = page_name.split(".")
        if spl[-1] not in [ "png", "jpg", "jpeg", "tiff"]:
            raise Exception("Invalid file extension")

        self.minio.client.upload_fileobj(io.BytesIO(batch), self.minio.bucket_name, str(uid)+".png")
        page_meta = schemas.responses.PageMetadata(id=uid, name=page_name, size=len(batch),
                                                   number=page_name.lstrip("0").split(".")[0], document_id=document_id,
                                                   attributes=None)
        await self.pg.save_page(page_meta)
        return uid

    async def delete_document(self, document_id:uuid.UUID)-> None:
        self.minio.client.delete_object(Bucket=self.minio.bucket_name, Key=str(document_id)+".png")

