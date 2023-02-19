NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
SERVICE_NAME=sm

build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm

run:
	make build
	go run main.go

test:
	@echo "$(OK_COLOR)==> Running tests$(NO_COLOR)"
	@go test -failfast -race -covermode=atomic -coverprofile=coverage.out ./...

test-short:
	@echo "$(OK_COLOR)==> Running short tests$(NO_COLOR)"
	@go test -short -failfast -race -covermode=atomic -coverprofile=coverage.out ./...

db:
	@PGPASSWORD=postgres psql -U postgres -d postgres --port 5430 --host localhost

migrate:
	@migrate -path="migrations" -database="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable" up

migrate-down:
	@migrate -path="migrations" -database="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable" down

seed:
	@go run seeds/seeder.go

reseed:
	@make migrate-down
	@make migrate
	@make seed

start:
	@echo "$(OK_COLOR)==> Bringing containers up for $(SERVICE_NAME)... $(NO_COLOR)"
	@docker-compose -f ./docker-compose.yml up -d

stop:
	@echo "$(OK_COLOR)==> Bringing containers down for $(SERVICE_NAME)... $(NO_COLOR)"
	@docker-compose -f ./docker-compose.yml down

reload: stop start

destroy:
	@echo "$(OK_COLOR)==> Bringing containers down and removing volumes for $(SERVICE_NAME)... $(NO_COLOR)"
	@docker-compose -f ./docker-compose.yml down --rmi all --volumes

ps:
	@echo "$(OK_COLOR)==> Checking containers status of $(SERVICE_NAME)... $(NO_COLOR)"
	@docker-compose -f ./docker-compose.yml ps
