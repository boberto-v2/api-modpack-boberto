# //TODO: Explain to Daniel how to use makefile and how we will implement the CI/CD using Makefile at server level


RED=\033[0;31m
NC=\033[0m
CYAN=\033[1;36m
user := root
database := test

dev:
	docker-compose up -d
	@echo "docker development setup started."
	
dev-build:
	docker-compose up --build
	@echo "docker compose image rebuilded."

stop:
	@echo "Stopping and clear all"
	docker-compose down
	@echo "Docker compose Stopped"

swagger:
	docker-compose exec -i app swag init
	@echo "Swagger doc generated"

test:
	@echo "${CYAN} Running API-MODPACK-TEST ${CYAN}"
	@echo "${RED}==> Running tests using docker-compose ${RED}"
	@docker-compose up -d
	@sleep 3 && \
		PG_URI="postgres://test:test@`docker-compose port postgres 5432`/test?sslmode=disable" \
		go test ./test_database -timeout 60s -cover -coverprofile=test_database/coverage.txt ./...
	@echo "${RED} coverate.txt is ready ${RED}"
	@echo "${RED} shutdown temporary database ${RED}"
	@docker-compose down