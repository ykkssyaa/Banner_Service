# Makefile

docker.start.components:
	docker compose up -d

docker.rebuild.components:
	docker-compose up -d --build app

# shutting down docker components
docker.stop:
	docker-compose down

export TEST_REDIS_URI=localhost:6376
export TEST_DB_URI=postgres://yks:yksadm@localhost:5434/postgres?sslmode=disable

test.integration:
	docker run --rm -d -p 5434:5432 --name test_db -e POSTGRES_PASSWORD=yksadm -e POSTGRES_USER=yks postgres:16
	docker run -d --rm -p 6376:6379 --name test_redis redis:7.2

	timeout 5

	migrate -path ./migrations -database "postgres://yks:yksadm@localhost:5434/postgres?sslmode=disable" up

	go test -v ./tests/

	docker stop test_db
	docker stop test_redis