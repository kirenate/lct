from dataclasses import dataclass

from repositories.pg_repository import PgRepository
from repositories.redis_repository import RedisRepository
from services.heath_service import HealthService
from presentations.web.presentation import Presentation


@dataclass
class Container:
    heath_service: HealthService
    redis_repository: RedisRepository

    @classmethod
    async def build_from_settings(cls) -> "Container":
        redis_repository = RedisRepository()
        pg_repository = PgRepository()
        heath_service = HealthService(redis_repository=redis_repository, pg_repository=pg_repository)

        return cls(
            heath_service=heath_service,
            redis_repository=redis_repository,
        )
