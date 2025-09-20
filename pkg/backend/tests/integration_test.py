from shared.containers import init_combat_container

container = init_combat_container()


class TestIntegration:
    async def test_health_check_ok(self) -> None:
        await container.heath_service.check()
