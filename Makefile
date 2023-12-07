SHELL := /bin/bash # Use bash syntax
include .envrc

## HELPERS

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/n]' && read ans && [ $${ans:-n} = y ]

## BUILD

current_time = $(shell date  +"%Y-%m-%dT%H:%M:%SZ")
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api



.PHONY: build/api/linux
build/api/linux:
	@echo 'Building cmd/api for linux'
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api

## DEVELOPMENT

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN}

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
	migrate -path ./migrations -database=${GREENLIGHT_DB_DSN} up

## QUALITY CONTROL
## audit: tidy dependencies and format, and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying dependencies'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor