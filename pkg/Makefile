update:
	git pull
	docker compose down
	docker compose up -d --build
	docker compose logs -f


backend-docker-exec:
	docker compose exec backend_web /bin/bash

infra-run:
	docker compose -f docker-compose-infra.yaml up -d