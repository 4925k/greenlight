SHELL := /bin/bash # Use bash syntax

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/n]' && read ans && [ $${ans:-n} = y ]

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api

## db/up helps to start the local db
/PHONY: db/up
db/up:
	/opt/homebrew/opt/postgresql@14/bin/postgres -D /opt/homebrew/var/postgresql@14

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

# db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migration'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up