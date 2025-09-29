from fastapi.responses import JSONResponse
from schemas.base import CamelizedBaseModel
import uuid
import datetime


class Attribute(CamelizedBaseModel):
    data:bytes                            #TODO : set correct json here

class PageMetadata(CamelizedBaseModel):
    id: uuid.UUID
    document_id: uuid.UUID
    name: str
    size: int
    number: int
    attributes: list[Attribute] | None

class Page(CamelizedBaseModel):
    id:uuid.UUID
    meta:PageMetadata

class DocumentMetadata(CamelizedBaseModel):
    id: uuid.UUID
    name: str
    size: int
    page_size: int
    pages: list[PageMetadata]
    created_at: datetime.datetime
    updated_at: datetime.datetime

class DocumentResponseSchema(CamelizedBaseModel):
    data: list[Page]
    meta: DocumentMetadata

class DocumentResponseMeta(CamelizedBaseModel):
    documents: list[DocumentMetadata]
    has_next: bool

class PageResponseMeta(CamelizedBaseModel):
    pages: list[PageMetadata]
    has_next: bool

