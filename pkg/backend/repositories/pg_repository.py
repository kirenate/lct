from dataclasses import dataclass

from shared.settings import app_settings
from sqlalchemy import text
from sqlalchemy.dialects import postgresql
from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy.schema import CreateTable


@dataclass
class PgRepository:
    def __post_init__(self) -> None:
        pool_size = 10
        max_overflow = 10
        self._aengine = create_async_engine(
            f"postgresql+asyncpg://{app_settings.pg.username}:{app_settings.pg.password}@"
            f"{app_settings.pg.host}:{app_settings.pg.port}/{app_settings.pg.database}",
            pool_size=pool_size,
            max_overflow=max_overflow,
        )

    async def health(self) -> None:
        async with self._aengine.connect() as session:
            result = await session.execute(text("select 1"))
            one = result.fetchone()
            if one is not None and one[0] != 1:
                raise Exception('Should be 1 from "select 1"')

    @staticmethod
    def compile_table(table) -> str:  # noqa: ANN001
        return str(CreateTable(table.__table__).compile(dialect=postgresql.dialect()))
