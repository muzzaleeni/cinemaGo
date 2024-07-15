include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${CINEMAGO_DB_DSN}

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${CINEMAGO_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${CINEMAGO_DB_DSN} up

# ==================================================================================== #
# DATABASE OPERATIONS
# ==================================================================================== #

## db/docker/init: create a Docker container for PostgreSQL, restore if backup exists, otherwise run migrations
.PHONY: db/docker/init
db/docker/init:
	@echo 'Creating Docker container for PostgreSQL...'
	docker run -d --name muzzyaqow -p 5432:5432 -e POSTGRES_PASSWORD=db_password -e POSTGRES_USER=postgres -e POSTGRES_DB=cinemago postgres
	@sleep 5 # Wait for the database to initialize
	@if [ -f "backup.sql" ]; then \
		echo 'Restoring PostgreSQL database from backup.sql...'; \
		cat backup.sql | docker exec -i muzzyaqow psql -U postgres; \
	else \
		echo 'No backup found, running migrations...'; \
		make db/migrations/up; \
	fi

## db/docker/delete: dump the PostgreSQL database data and delete the container
.PHONY: db/docker/delete
db/docker/delete: confirm
	@echo 'Dumping PostgreSQL database data...'
	docker exec -t muzzyaqow pg_dumpall -c -U postgres > backup.sql
	@echo 'Deleting Docker container for PostgreSQL...'
	docker rm -f muzzyaqow

## db/dump: dump the PostgreSQL database data
.PHONY: db/dump
db/dump:
	@echo 'Dumping PostgreSQL database data...'
	docker exec -t muzzyaqow pg_dumpall -c -U postgres > dump_$$(date +%Y-%m-%d_%H_%M_%S).sql

## db/restore file=<path/to/dump.sql>: restore the PostgreSQL database data
.PHONY: db/restore
db/restore:
	@echo 'Restoring PostgreSQL database data from ${file}...'
	cat ${file} | docker exec -i muzzyaqow psql -U postgres

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
