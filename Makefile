PGSQL_LOCAL_CONN = "host=localhost user=postgres password=postgres port=5432 dbname=db sslmode=disable"
PGSQL_STAGING_CONN := ""
PGSQL_PRODUCTION_CONN := ""
GOOSE_LINK := github.com/pressly/goose/v3/cmd/goose@latest
GOOSE_TAGS := 'no_mssql no_redshift no_tidb no_clickhouse no_vertica no_mysql no_sqlite3 no_ydb'

help: ## Help target
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
build:
	go build -o ./bin/app ./cmd/http/http.go
run: build
	@env CONFIG_PATH="./internal/config/development.yaml" ./bin/app
watch:
	@env CONFIG_PATH="./internal/config/development.yaml" air
test: ## Run test suite
	@env CONFIG_PATH="./internal/config/test.yaml" go test -coverprofile=cover.out -count=1 -v -p 1 ./...
coverage: test ## Shows test coverage
	@go tool cover -html=cover.out
db: ## Start local database in a docker container
	@docker compose -f ./testing-database/docker-compose-pg.yaml up
stop_db: ## Stop local database in a docker container
	@docker compose -f ./testing-database/docker-compose-pg.yaml down

local_upmigrate: ## Run database migrations on the local database UP
	@go run -tags=$(GOOSE_TAGS) $(GOOSE_LINK) -dir "./migrations" postgres $(PGSQL_LOCAL_CONN) up
local_downmigrate: ## Run database migrations on the local database DOWN
	@go run -tags=$(GOOSE_TAGS) $(GOOSE_LINK) -dir "./migrations" postgres $(PGSQL_LOCAL_CONN) down
local_reset: ## Run database migrations on the local database DOWN-TO version 0
	@go run -tags=$(GOOSE_TAGS) $(GOOSE_LINK) -dir "./migrations" postgres $(PGSQL_LOCAL_CONN) down-to 0
local_status: ## Print migrations status on the local database
	@go run -tags=$(GOOSE_TAGS) $(GOOSE_LINK) -dir "./migrations" postgres $(PGSQL_LOCAL_CONN) status

staging_upmigrate: ## Run database migrations on the staging database UP
	@go run -tags=$(GOOSE_TAGS) $(GOOSE_LINK) -dir "./migrations" postgres $(PGSQL_STAGING_CONN) up
staging_downmigrate: ## Run database migrations on the staging database DOWN
	@go run -tags=$(GOOSE_TAGS) $(GOOSE_LINK) -dir "./migrations" postgres $(PGSQL_STAGING_CONN) down
staging_reset: ## Run database migrations on the staging database DOWN-TO version 0
	@go run -tags=$(GOOSE_TAGS) $(GOOSE_LINK) -dir "./migrations" postgres $(PGSQL_STAGING_CONN) down-to 0
staging_status: ## Print migrations status on the staging database
	@go run -tags=$(GOOSE_TAGS) $(GOOSE_LINK) -dir "./migrations" postgres $(PGSQL_STAGING_CONN) status

prod_upmigrate: ## Run database migrations on the production database UP
	@go run -tags=$(GOOSE_TAGS) $(GOOSE_LINK) -dir "./migrations" postgres $(PGSQL_PRODUCTION_CONN) up
prod_downmigrate: ## Run database migrations on the production database DOWN
	@go run -tags=$(GOOSE_TAGS) $(GOOSE_LINK) -dir "./migrations" postgres $(PGSQL_PRODUCTION_CONN) down
prod_reset: ## Run database migrations on the production database DOWN-TO version 0
	@go run -tags=$(GOOSE_TAGS) $(GOOSE_LINK) -dir "./migrations" postgres $(PGSQL_PRODUCTION_CONN) down-to 0
prod_status: ## Print migrations status on the production database
	@go run -tags=$(GOOSE_TAGS) $(GOOSE_LINK) -dir "./migrations" postgres $(PGSQL_PRODUCTION_CONN) status

