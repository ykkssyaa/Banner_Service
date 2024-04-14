# Makefile

docker.start.components:
	docker compose up -d

docker.rebuild.components:
	docker compose up -d --build app

docker.stop:
	docker compose down

migrate.up:
	migrate -path ./migrations -database "postgres://yks:yksadm@localhost:5432/postgres?sslmode=disable" up

migrate.down:
	migrate -path ./migrations -database "postgres://yks:yksadm@localhost:5432/postgres?sslmode=disable" down

export TEST_PORT=8090
export TEST_REDIS_URI=localhost:6376
export TEST_DB_URI=postgres://yks:yksadm@localhost:5434/postgres?sslmode=disable

test.integrations:
	docker compose -f docker-compose.test.yml up -d

	go test -v ./tests/ || (docker stop test_db && docker stop test_redis && exit 1)

	docker compose -f docker-compose.test.yml rm