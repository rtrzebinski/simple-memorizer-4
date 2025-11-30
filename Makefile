NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
SERVICE_NAME=sm
TIMEZONE=Europe/Warsaw
HOME ?= $(shell echo $$HOME)

MIGRATE_URL=postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable&timezone=$(TIMEZONE)
MIGRATE_CMD=migrate -path="migrations" -database="$(MIGRATE_URL)"

default: help

.PHONY: dev proto

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
	@brew install golang-migrate
	@brew install protobuf
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

dev: ## Prepare local dev environment
	@make stop
	@make start
	@make database & make pubsub & make keycloak & wait
	@make seed

start: ## Start docker-compose containers
	@docker-compose -f ./dev/docker-compose.yml up -d --remove-orphans

stop: ## Stop docker-compose containers
	@docker-compose -f ./dev/docker-compose.yml down --remove-orphans -t 0

database: start ## Migrate the database when it's ready
	@until docker exec dev-sm-db-1 pg_isready 2>&1 | grep -q "accepting connections"; do \
          sleep 1; \
          echo "[Database] Waiting for Postgres.."; \
         done
	@make migrate

keycloak: ## Configure Keycloak and create a dev user
	@chmod +x dev/keycloak/keycloak-init.sh
	@./dev/keycloak/keycloak-init.sh
	@echo "[Keycloak] Configured"

pubsub: ## Create PubSub topic and subscription
	@PUBSUB_EMULATOR_HOST=0.0.0.0:8088 go run ./dev/pubsub/pubsub.go
	@echo "[PubSub] Configured"

ps: ## Show running containers
	@docker-compose -f ./dev/docker-compose.yml ps

restart: stop start ## Stop and start containers

destroy: ## Stop containers and remove volumes
	@docker-compose -f ./dev/docker-compose.yml down --rmi all --volumes

migrate: ## Run db migrations (migrate up)
	@$(MIGRATE_CMD) up
	@echo "[Database] Migrated"

migrate-down: ## Revert db migrations (migrate down)
	@$(MIGRATE_CMD) down

migrate-drop: ## Drop db without confirmation (migrate drop)
	@$(MIGRATE_CMD) drop -f

seed: ## Seed the database with some example data
	@go run dev/seeder/seeder.go
	@echo "[Database] Seeded"

reseed: ## Destroy, recreate and seed the database (no confirmation)
	@make migrate-down
	@make migrate
	@make seed

db: ## Db CLI client connection
	@PGPASSWORD=postgres PGTZ="$(TIMEZONE)" psql -U postgres -d postgres --port 5430 --host localhost

build-web: ## Build web (client and server)
	@GOARCH=wasm GOOS=js go build -o web/app.wasm github.com/rtrzebinski/simple-memorizer-4/cmd/web
	@go build -o bin/sm4-web github.com/rtrzebinski/simple-memorizer-4/cmd/web
	@date > version

run-web: ## Build and run web locally
	@make build-web
	@PUBSUB_EMULATOR_HOST=0.0.0.0:8088 ./bin/sm4-web

build-worker: ## Build worker
	@date > version
	@PUBSUB_EMULATOR_HOST=0.0.0.0:8088 go build -o bin/sm4-worker github.com/rtrzebinski/simple-memorizer-4/cmd/worker

run-worker: ## Run worker locally
	@make build-worker
	@PUBSUB_EMULATOR_HOST=0.0.0.0:8088 ./bin/sm4-worker

build-auth: ## Build auth
	@date > version
	@go build -o bin/sm4-auth github.com/rtrzebinski/simple-memorizer-4/cmd/auth

run-auth: ## Run auth locally
	@make build-auth
	@PUBSUB_EMULATOR_HOST=0.0.0.0:8088 ./bin/sm4-auth

run: ## Build and run all services locally
	@make run-web & make run-worker & make run-auth & echo "$(OK_COLOR)==> Running on http://localhost:8000 $(NO_COLOR)" & wait

logs-keycloak: ## Show keycloak container logs
	docker-compose -f dev/docker-compose.yml logs sm-keycloak

proto: ## Generate protobuf files
	@protoc --go_out=./generated --go_opt=paths=source_relative proto/events/*.proto
	@protoc --go_out=./generated --go-grpc_out=./generated --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/grpc/*.proto

test: ## Test all (unit + integration)
	@go test -failfast -race -covermode=atomic -coverprofile=coverage.out -ldflags=-extldflags=-Wl,-ld_classic ./...

test-short: ## Test short (unit)
	@go test -short -failfast -race -covermode=atomic -coverprofile=coverage.out -ldflags=-extldflags=-Wl,-ld_classic ./...

go-deps-update: ## Update GO dependencies
	@go get -u ./...
	@go vet ./...

k8s-start: ## Kubernetes create all objects (Docker hub tag 'latest' image)
	@kubectx docker-desktop
	@kubectl create namespace sm4
	@mkdir -p $(HOME)/sm4-db
	@kubectl -n sm4 apply -f k8s/pubsub-deployment.yaml
	@kubectl -n sm4 apply -f k8s/web-deployment.yaml
	@kubectl -n sm4 apply -f k8s/worker-deployment.yaml
	@kubectl -n sm4 apply -f k8s/auth-deployment.yaml
	@envsubst < k8s/db-deployment.yaml | kubectl -n sm4 apply -f -
	@kubectl -n sm4 apply -f k8s/keycloak-deployment.yaml
	@kubectl -n sm4 apply -f k8s/db-migration-job.yaml
	@mkdir -p $(HOME)/sm4-db-backup
	@envsubst < k8s/db-backup-cronjob.yaml  | kubectl -n sm4 apply -f -
	@kubectl apply -f k8s/metrics-server.yaml
	@echo "==> Keycloak UI http://localhost:30002 (admin:change_me)"
	@echo "$(OK_COLOR)==> Running on http://localhost:9000 $(NO_COLOR)"

k8s-status: ## Kubernetes show all objects
	@kubectl -n sm4 get all

k8s-stop: ## Kubernetes delete all objects
	@kubectl delete namespace sm4 --ignore-not-found=true
	@kubectl delete -f k8s/metrics-server.yaml --ignore-not-found=true

k8s-reset: ## Kubernetes stop and start
	@make k8s-stop
	@make k8s-start

k8s-rollout: ## Kubernetes rollout (Docker hub tag 'latest' image)
	@kubectl -n sm4 rollout restart deployment.apps/sm4-auth
	@kubectl -n sm4 rollout restart deployment.apps/sm4-web
	@kubectl -n sm4 rollout restart deployment.apps/sm4-worker
	@echo "$(OK_COLOR)==> Running on http://localhost:9000 $(NO_COLOR)"

k8s-logs-auth: ## Kubernetes auth logs
	@kubectl -n sm4 logs -l app=sm4-auth -f

k8s-logs-web: ## Kubernetes web logs
	@kubectl -n sm4 logs -l app=sm4-web -f

k8s-logs-worker: ## Kubernetes worker logs
	@kubectl -n sm4 logs -l app=sm4-worker -f

k8s-logs: ## Kubernetes logs
	@make k8s-logs-auth & make k8s-logs-web & make k8s-logs-worker & wait

k8s-sh: ## Kubernetes web app shell
	@kubectl -n sm4 exec -it deployment.apps/sm4-web -- sh

k8s-db: ## Kubernetes db cli
	@PGPASSWORD=postgres psql -U postgres -d postgres --port 30001 --host localhost

k8s-db-migrate: ## Kubernetes db migrate
	@kubectl -n sm4 apply -f k8s/db-migration-job.yaml
