GO ?= GO111MODULE=on CGO_ENABLED=1 go

run:
	$(GO) run backend/main.go

build:
	$(GO) build -o=bootstrap backend/main.go

test:
	$(GO) test ./...

update:
	git pull
	docker compose down
	docker compose up -d --build
	docker compose logs -f


backend-docker-exec:
	docker compose exec backend_web /bin/bash

infra-run:
	docker compose -f docker-compose-infra.yaml up -d
