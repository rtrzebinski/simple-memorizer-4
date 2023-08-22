NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
SERVICE_NAME=sm
HOME ?= $(shell echo $$HOME)

default: help

.PHONY: all dev

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

dev: ## Prepare local dev environment (stop + start + migrate + seed)
	@echo "$(OK_COLOR)==> Prepare dev environment for $(SERVICE_NAME)... $(NO_COLOR)"
	@make stop
	@make start
	@echo "$(OK_COLOR)==> Waiting for the db to be ready... $(NO_COLOR)"
	@sleep 1
	@make migrate
	@make seed
	@echo "$(OK_COLOR)==> Completed $(NO_COLOR)"

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
	@migrate -path="migrations" -database="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable" up

migrate-down: ## Revert db migrations (migrate down)
	@echo "$(OK_COLOR)==> Reverting db migrations for $(SERVICE_NAME)... $(NO_COLOR)"
	@migrate -path="migrations" -database="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable" down

migrate-drop: ## Drop db without confirmation (migrate drop)
	@echo "$(OK_COLOR)==> Dropping db migrations for $(SERVICE_NAME)... $(NO_COLOR)"
	@migrate -path="migrations" -database="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable" drop -f

seed: ## Seed the database with some example data
	@echo "$(OK_COLOR)==> Seeding the db for $(SERVICE_NAME)... $(NO_COLOR)"
	@go run seeds/seeder.go

reseed: ## Destroy, recreate and seed the database (no confirmation)
	@make migrate-down
	@make migrate
	@make seed

db: ## Db CLI client connection
	@echo "$(OK_COLOR)==> Connecting to the db of $(SERVICE_NAME)... $(NO_COLOR)"
	@PGPASSWORD=postgres psql -U postgres -d postgres --port 5430 --host localhost

build: ## Build client and server
	@echo "$(OK_COLOR)==> Building client and server for $(SERVICE_NAME)... $(NO_COLOR)"
	@GOARCH=wasm GOOS=js go build -o web/app.wasm github.com/rtrzebinski/simple-memorizer-4/cmd/simple-memorizer
	@go build -o bin/simple-memorizer github.com/rtrzebinski/simple-memorizer-4/cmd/simple-memorizer
	@date > version
	@echo "$(OK_COLOR)==> Completed $(NO_COLOR)"

run: ## Build and run locally
	@make build
	@echo "$(OK_COLOR)==> Running on http://localhost:8000 $(NO_COLOR)"
	@go run cmd/simple-memorizer/main.go

test: ## Test all
	@echo "$(OK_COLOR)==> Running tests for $(SERVICE_NAME)... $(NO_COLOR)"
	@go test -failfast -race -covermode=atomic -coverprofile=coverage.out ./...
	@echo "$(OK_COLOR)==> Completed $(NO_COLOR)"

test-short: ## Test short (unit)
	@echo "$(OK_COLOR)==> Running short tests for $(SERVICE_NAME)... $(NO_COLOR)"
	@go test -short -failfast -race -covermode=atomic -coverprofile=coverage.out ./...
	@echo "$(OK_COLOR)==> Completed $(NO_COLOR)"

k8s-deploy-all: ## Kubernetes deploy all objects
	@make k8s-deploy
	@make k8s-db-backup-cronjob-deploy

k8s-delete-all: ## Kubernetes delete all objects
	@make k8s-delete
	@make k8s-db-backup-cronjob-delete

k8s-deploy: ## Kubernetes deploy
	@mkdir -p $(HOME)/sm4-db
	@kubectl apply -f k8s/local-web-deployment.yaml
	@envsubst < k8s/local-db-deployment.yaml | kubectl apply -f -
	@kubectl apply -f k8s/local-db-migration-job.yaml
	@echo "$(OK_COLOR)==> Running on http://localhost:9000 $(NO_COLOR)"

k8s-rollout: ## Kubernetes rollout
	kubectl rollout restart deployment.apps/sm4-web
	@echo "$(OK_COLOR)==> Running on http://localhost:9000 $(NO_COLOR)"

k8s-delete: ## Kubernetes delete
	@kubectl delete -f k8s/local-web-deployment.yaml --ignore-not-found=true
	@kubectl delete -f k8s/local-db-deployment.yaml --ignore-not-found=true

k8s-logs: ## Kubernetes web app logs
	@kubectl logs -l app=sm4-web -f

k8s-sh: ## Kubernetes web app shell
	@kubectl exec -it deployment.apps/sm4-web -- sh

k8s-db: ## Kubernetes db cli
	@PGPASSWORD=postgres psql -U postgres -d postgres --port 30001 --host localhost

k8s-db-migrate: ## Kubernetes db migrate
	@kubectl apply -f k8s/local-db-migration-job.yaml

k8s-db-seed: ## Kubernetes db seed
	@kubectl exec deployment.apps/sm4-web -- make seed

k8s-db-backup-cronjob-deploy: ## Kubernetes db backup CRON job deploy
	@mkdir -p $(HOME)/sm4-db-backup
	@envsubst < k8s/local-db-backup-cronjob.yaml  | kubectl apply -f -

k8s-db-backup-cronjob-delete: ## Kubernetes db backup CRON job delete
	@kubectl delete -f k8s/local-db-backup-cronjob.yaml --ignore-not-found=true
