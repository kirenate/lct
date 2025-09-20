import asyncio

import uvicorn
from presentations.web.app import create_app
from shared.containers import Container
from shared.settings import app_settings


async def main() -> None:
    container = await Container.build_from_settings()
    app = await create_app(container)

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
