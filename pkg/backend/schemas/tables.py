import sqlalchemy
from fastapi.responses import JSONResponse
from schemas.base import CamelizedBaseModel
import uuid
import datetime
from persistent.database import Base
from schemas.responses import Attribute, PageMetadata


class PageMetadataTable(Base):
    table_name:str="page_metadata"
    id: uuid.UUID
    document_id: uuid.UUID
    name: str
    size: int
    number: int
    attributes: list[Attribute] | None


class DocumentMetadataTable(Base):
    table_name:str="document_metadata"
    id: uuid.UUID
    name: str
    size: int
    page_size: int
    pages: list[PageMetadata]
    created_at: datetime.datetime
    updated_at: datetime.datetime

