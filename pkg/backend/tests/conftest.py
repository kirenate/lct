import pytest
from shared.containers import Container, init_combat_container


@pytest.fixture()
async def combat_container() -> Container:
    return await init_combat_container()
