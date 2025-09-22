import asyncio

import uvicorn
from presentations.web.app import create_app
from shared.containers import Container
from shared.settings import app_settings
from services.batch_processing_service import Service


async def main() -> None:
    container = await Container.build_from_settings()
    service = Service()
    app = await create_app(container, service)

    server = uvicorn.Server(
        uvicorn.Config(
            app,
            host=app_settings.uvicorn.host,
            port=app_settings.uvicorn.port,
            workers=app_settings.uvicorn.workers,
            ssl_keyfile=app_settings.uvicorn.ssl_keyfile,
            ssl_certfile=app_settings.uvicorn.ssl_certfile,
        )
    )
    await server.serve()


if __name__ == "__main__":
    asyncio.run(main())
