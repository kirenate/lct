from dataclasses import dataclass

from repository.pg_repository import PgRepository
from repository.redis_repository import RedisRepository
from service.heath_service import HeathService


@dataclass
class Container:
    heath_service: HeathService
    redis_repository: RedisRepository

    @classmethod
    async def build_from_settings(cls) -> "Container":
        redis_repository = RedisRepository()
        pg_repository = PgRepository()
        heath_service = HeathService(redis_repository=redis_repository, pg_repository=pg_repository)

        return cls(
            heath_service=heath_service,
            redis_repository=redis_repository,
        )
