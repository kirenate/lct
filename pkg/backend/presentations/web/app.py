from fastapi import FastAPI, HTTPException, status, UploadFile
from fastapi.middleware.cors import CORSMiddleware
from fastapi_cache import FastAPICache
from fastapi_cache.backends.redis import RedisBackend
from shared.containers import Container
from presentations.web.presentation import Presentation

async def create_app(container: Container, presentation: Presentation) -> FastAPI:
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
        FastAPICache.init(RedisBackend(presentation.service.redis.ar), prefix="fastapi-cache")

    @app.get("/health")
    async def check_server_health() -> bool:
        try:
            await container.heath_service.check()
        except Exception as exc:
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="{exc.__class__.__name__}: {str(exc)}"
            ) from exc

        return True

    @app.post("/batches")
    async def upload_batch(batch: UploadFile) -> bool:
        try:
            await presentation.upload_batch(batch)
        except Exception as e:
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="{exc.__class__.__name__}: {str(e)}"
            ) from e

        return True

    @app.delete("/batches")
    async def delete_batch(batch_name:str) ->bool:
        try:
            await presentation.delete_batch(batch_name)
        except Exception as e:
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="{exc.__class__.__name__}: {str(e)}"
            ) from e
        return True

    return app
