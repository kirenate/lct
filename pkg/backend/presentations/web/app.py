from fastapi import FastAPI, HTTPException, status, UploadFile, Query, Body
from fastapi.middleware.cors import CORSMiddleware
from fastapi_cache import FastAPICache
from fastapi_cache.backends.redis import RedisBackend
from shared.containers import Container
from services.filter_service import Filter
import uuid
from schemas.responses import PageResponseMeta, DocumentResponseMeta
from services.batch_processing_service import Service

async def create_app(container: Container, service: Service) -> FastAPI:
    app = FastAPI(title="MISIS_MCs", root_path="/api")
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    @app.on_event("startup")
    async def startup() -> None:
        FastAPICache.init(RedisBackend(service.redis.ar), prefix="fastapi-cache")

    @app.get("/health")
    async def check_server_health() -> bool:
        try:
            await container.heath_service.check()
        except Exception as exc:
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="{exc.__class__.__name__}: {str(exc)}"
            ) from exc

        return True

    @app.post("/documents")
    async def upload_document(file: UploadFile) -> uuid.UUID:
        res = await service.save_document(bytes(await file.read()), str(file.filename))
        await file.seek(0)
        return res


    @app.delete("/documents/{documentId: uuid.UUID }")
    async def delete_document(document_id: uuid.UUID) -> None:
        await service.delete_document(document_id)


    @app.post("/documents/get")
    async def get_documents(page:int = Query(...), page_size: int = Query(..., alias="pageSize"),
                            sort_by: str = Query(..., alias="sortBy"),
                            filters: list[Filter] = Body(...))->DocumentResponseMeta:
        return DocumentResponseMeta()

    @app.post("/documents/{documentId: uuid.UUID}/pages/get")
    async def get_document_pages(document_id: uuid.UUID = Query(...,alias="documentId"), page:int = Query(...),
                                 page_size:int = Query(...,alias="pageSize"),
                                 sort_by:str = Query(..., alias="sortBy"),
                                 filters:list[Filter] = Body(...)) ->PageResponseMeta:
        resp = PageResponseMeta()
        return resp


    return app
