import asyncio

import uvicorn
from more_itertools.more import bucket

from presentations.web.app import create_app
from shared.containers import Container
from shared.settings import app_settings
from presentations.web.presentation import Presentation
from services.batch_processing_service import Service


async def main() -> None:
    container = await Container.build_from_settings()
    service = Service()
    presentation = Presentation(service)
    app = await create_app(container, presentation)

    server = uvicorn.Server(
        uvicorn.Config(
            app,
            host=app_settings.uvicorn.host,
            port=app_settings.uvicorn.port,
            workers=app_settings.uvicorn.workers,
        )
    )
    await server.serve()


if __name__ == "__main__":
    asyncio.run(main())
