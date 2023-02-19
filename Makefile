NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
SERVICE_NAME=sm

default: help

help: ## Show this help
	@IFS=$$'\n' ; \
    help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##/:/'`); \
    printf "%-30s %s\n" "target" "help" ; \
    printf "%-30s %s\n" "------" "----" ; \
    for help_line in $${help_lines[@]}; do \
        IFS=$$':' ; \
        help_split=($$help_line) ; \
        help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
        help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
        printf '\033[36m'; \
        printf "%-30s %s" $$help_command ; \
        printf '\033[0m'; \
        printf "%s\n" $$help_info; \
    done

build: ## Build client and server
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	go build

run: ## Build and run locally
	make build
	go run main.go

test: ## Test all
	@echo "$(OK_COLOR)==> Running tests$(NO_COLOR)"
	@go test -failfast -race -covermode=atomic -coverprofile=coverage.out ./...

test-short: ## Test short (unit)
	@echo "$(OK_COLOR)==> Running short tests$(NO_COLOR)"
	@go test -short -failfast -race -covermode=atomic -coverprofile=coverage.out ./...

db: ## Database CLI client connection
	@PGPASSWORD=postgres psql -U postgres -d postgres --port 5430 --host localhost

migrate: ## Run migrations (migrate up)
	@migrate -path="migrations" -database="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable" up

migrate-down: ## Revert migrations (migrate down)
	@migrate -path="migrations" -database="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable" down

seed: ## Seed the database with example data
	@go run seeds/seeder.go

reseed: ## Destroy, recreate and seed database
	@make migrate-down
	@make migrate
	@make seed

start: ## Start containers (docker compose up)
	@echo "$(OK_COLOR)==> Bringing containers up for $(SERVICE_NAME)... $(NO_COLOR)"
	@docker-compose -f ./docker-compose.yml up -d

stop: ## Stop containers (docker compose down)
	@echo "$(OK_COLOR)==> Bringing containers down for $(SERVICE_NAME)... $(NO_COLOR)"
	@docker-compose -f ./docker-compose.yml down

reload: stop start ## Stop and start again

destroy: ## Stop and remove volumes
	@echo "$(OK_COLOR)==> Bringing containers down and removing volumes for $(SERVICE_NAME)... $(NO_COLOR)"
	@docker-compose -f ./docker-compose.yml down --rmi all --volumes

ps: ## Show running containers
	@echo "$(OK_COLOR)==> Checking containers status of $(SERVICE_NAME)... $(NO_COLOR)"
	@docker-compose -f ./docker-compose.yml ps
