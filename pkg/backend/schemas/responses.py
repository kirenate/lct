from fastapi.responses import JSONResponse
from schemas.base import CamelizedBaseModel



class Attribute(CamelizedBaseModel):
    data:str                            #TODO : set correct json here

class Page(CamelizedBaseModel):
    id: uuid.UUID
    name: str
    size: integer
    number: integer
    attributes: list[Attribute]

class Document(CamelizedBaseModel):
    id: uuid.UUID
    name: str
    size: integer
    created_at: datetime
    updated_at: datetime

class DocumentResponseSchema(CamelizedBaseModel):
    data: bytes
    meta: Document

class DocumentResponse(CamelizedBaseModel):
    document: list[DocumentResponseSchema]



