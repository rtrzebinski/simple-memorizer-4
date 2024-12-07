NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
SERVICE_NAME=sm
TIMEZONE=Europe/Warsaw
HOME ?= $(shell echo $$HOME)

default: help

.PHONY: all dev proto

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

deps: ## Install local environment dependencies
	@echo "$(OK_COLOR)==> Installing local environment dependencies for $(SERVICE_NAME)... $(NO_COLOR)"
	@brew install golang-migrate
	@brew install protobuf
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

dev: ## Prepare local dev environment (stop + start + migrate + seed)
	@echo "$(OK_COLOR)==> Preparing dev environment for $(SERVICE_NAME)... $(NO_COLOR)"
	@make stop
	@make start
	@echo "$(OK_COLOR)==> Waiting for the db to be ready... $(NO_COLOR)"
	@until docker exec dev-sm-db-1 pg_isready 2>&1 | grep -q "accepting connections"; do \
          sleep 1; \
          echo "still waiting..."; \
         done
	@make migrate
	@make seed
	@make pubsub
	@echo "$(OK_COLOR)==> Completed $(NO_COLOR)"

pubsub: ## Create PubSub topic and subscription
	@echo "$(OK_COLOR)==> Setting up PubSub...$(NO_COLOR)"
	@PUBSUB_EMULATOR_HOST=0.0.0.0:8088 go run ./dev/pubsub/pubsub.go

start: ## Start docker-compose containers
	@echo "$(OK_COLOR)==> Bringing containers up for $(SERVICE_NAME)... $(NO_COLOR)"
	@docker-compose -f ./dev/docker-compose.yml up -d --remove-orphans

stop: ## Stop docker-compose containers
	@echo "$(OK_COLOR)==> Bringing containers down for $(SERVICE_NAME)... $(NO_COLOR)"
	@docker-compose -f ./dev/docker-compose.yml down --remove-orphans

ps: ## Show running containers
	@echo "$(OK_COLOR)==> Checking containers status of $(SERVICE_NAME)... $(NO_COLOR)"
	@docker-compose -f ./dev/docker-compose.yml ps

restart: stop start ## Stop and start containers

destroy: ## Stop containers and remove volumes
	@echo "$(OK_COLOR)==> Bringing containers down and removing volumes for $(SERVICE_NAME)... $(NO_COLOR)"
	@docker-compose -f ./dev/docker-compose.yml down --rmi all --volumes

migrate: ## Run db migrations (migrate up)
	@echo "$(OK_COLOR)==> Running db migrations for $(SERVICE_NAME)... $(NO_COLOR)"
	@migrate -path="migrations" -database="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable&timezone=$(TIMEZONE)" up

migrate-down: ## Revert db migrations (migrate down)
	@echo "$(OK_COLOR)==> Reverting db migrations for $(SERVICE_NAME)... $(NO_COLOR)"
	@migrate -path="migrations" -database="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable&timezone=$(TIMEZONE)" down

migrate-drop: ## Drop db without confirmation (migrate drop)
	@echo "$(OK_COLOR)==> Dropping db migrations for $(SERVICE_NAME)... $(NO_COLOR)"
	@migrate -path="migrations" -database="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable&timezone=$(TIMEZONE)" drop -f

seed: ## Seed the database with some example data
	@echo "$(OK_COLOR)==> Seeding the db for $(SERVICE_NAME)... $(NO_COLOR)"
	@go run dev/seeder/seeder.go

reseed: ## Destroy, recreate and seed the database (no confirmation)
	@make migrate-down
	@make migrate
	@make seed

db: ## Db CLI client connection
	@echo "$(OK_COLOR)==> Connecting to the db of $(SERVICE_NAME)... $(NO_COLOR)"
	@PGPASSWORD=postgres PGTZ="$(TIMEZONE)" psql -U postgres -d postgres --port 5430 --host localhost

build-web: ## Build client and server
	@echo "$(OK_COLOR)==> Building client and server for $(SERVICE_NAME)... $(NO_COLOR)"
	@GOARCH=wasm GOOS=js go build -o web/app.wasm github.com/rtrzebinski/simple-memorizer-4/cmd/web
	@go build -o bin/sm4-web github.com/rtrzebinski/simple-memorizer-4/cmd/web
	@date > version
	@echo "$(OK_COLOR)==> Completed $(NO_COLOR)"

run-web: ## Build and run web locally
	@make build-web
	@echo "$(OK_COLOR)==> Running on https://localhost:8000 $(NO_COLOR)"
	@PUBSUB_EMULATOR_HOST=0.0.0.0:8088 go run cmd/web/main.go

build-worker: ## Build and run worker locally
	@PUBSUB_EMULATOR_HOST=0.0.0.0:8088 go build -o bin/sm4-worker github.com/rtrzebinski/simple-memorizer-4/cmd/worker

run-worker: ## Build and run worker locally
	@PUBSUB_EMULATOR_HOST=0.0.0.0:8088 go run cmd/worker/main.go

run: ## Build and run all services locally
	@make build-web
	@echo "$(OK_COLOR)==> Running on https://localhost:8000 $(NO_COLOR)"
	@PUBSUB_EMULATOR_HOST=0.0.0.0:8088 go run cmd/web/main.go & PUBSUB_EMULATOR_HOST=0.0.0.0:8088 go run cmd/worker/main.go & wait

proto: ## Generate protobuf files
	@echo "$(OK_COLOR)==> Generating protobuf files for $(SERVICE_NAME)... $(NO_COLOR)"
	@protoc --go_out=./generated --go_opt=paths=source_relative proto/events/*.proto

test: ## Test all
	@echo "$(OK_COLOR)==> Running tests for $(SERVICE_NAME)... $(NO_COLOR)"
	@go test -failfast -race -covermode=atomic -coverprofile=coverage.out -ldflags=-extldflags=-Wl,-ld_classic ./...
	@echo "$(OK_COLOR)==> Completed $(NO_COLOR)"

test-short: ## Test short (unit)
	@echo "$(OK_COLOR)==> Running short tests for $(SERVICE_NAME)... $(NO_COLOR)"
	@go test -short -failfast -race -covermode=atomic -coverprofile=coverage.out -ldflags=-extldflags=-Wl,-ld_classic ./...
	@echo "$(OK_COLOR)==> Completed $(NO_COLOR)"

go-deps-update: ## Update GO dependencies
	@echo "$(OK_COLOR)==> Updating GO dependencies $(SERVICE_NAME)... $(NO_COLOR)"
	@go get -u ./...
	@go vet ./...

k8s-start: ## Kubernetes create all objects (Docker hub tag 'latest' image)
	@kubectx docker-desktop
	@kubectl create namespace sm4
	@mkdir -p $(HOME)/sm4-db
	@kubectl -n sm4 apply -f k8s/pubsub-deployment.yaml
	@kubectl -n sm4 apply -f k8s/web-deployment.yaml
	@kubectl -n sm4 apply -f k8s/worker-deployment.yaml
	@envsubst < k8s/db-deployment.yaml | kubectl -n sm4 apply -f -
	@kubectl -n sm4 apply -f k8s/db-migration-job.yaml
	@mkdir -p $(HOME)/sm4-db-backup
	@envsubst < k8s/db-backup-cronjob.yaml  | kubectl -n sm4 apply -f -
	@kubectl apply -f k8s/metrics-server.yaml
	@echo "$(OK_COLOR)==> Running on https://localhost:9000 $(NO_COLOR)"

k8s-status: ## Kubernetes show all objects
	@kubectl -n sm4 get all

k8s-stop: ## Kubernetes delete all objects
	@kubectl delete namespace sm4 --ignore-not-found=true
	@kubectl delete -f k8s/metrics-server.yaml --ignore-not-found=true

k8s-reset: ## Kubernetes stop and start
	@make k8s-stop
	@make k8s-start

k8s-rollout: ## Kubernetes rollout (Docker hub tag 'latest' image)
	@kubectl -n sm4 rollout restart deployment.apps/sm4-web
	@kubectl -n sm4 rollout restart deployment.apps/sm4-worker
	@echo "$(OK_COLOR)==> Running on http://localhost:9000 $(NO_COLOR)"

k8s-logs-web: ## Kubernetes web logs
	@kubectl -n sm4 logs -l app=sm4-web -f

k8s-logs-worker: ## Kubernetes worker logs
	@kubectl -n sm4 logs -l app=sm4-worker -f

k8s-logs: ## Kubernetes logs
	@make k8s-logs-worker & make k8s-logs-web & wait

k8s-sh: ## Kubernetes web app shell
	@kubectl -n sm4 exec -it deployment.apps/sm4-web -- sh

k8s-db: ## Kubernetes db cli
	@PGPASSWORD=postgres psql -U postgres -d postgres --port 30001 --host localhost

k8s-db-migrate: ## Kubernetes db migrate
	@kubectl -n sm4 apply -f k8s/db-migration-job.yaml
