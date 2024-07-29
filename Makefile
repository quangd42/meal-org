# Change these variables as necessary.
MAIN_PACKAGE_PATH := .
BINARY_NAME := planner_server
SCHEMA_PATH := sql/schema
PG_DATABASE_URL := postgres://quang-dang:@localhost:5432/postgres?sslmode=disable
ifneq (,$(wildcard ./.env))
		include .env
		export
endif

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run github.com/securego/gosec/v2/cmd/gosec@latest -exclude-generated -exclude-dir=scripts ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	## go test -race -buildvcs -vet=off ./...


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## sqlc: generate database code with sqlc
.PHONY: sqlc
sqlc:
	sqlc generate

## db/drop: drop local db
.PHONY: db/drop
db/drop:
	psql ${PG_DATABASE_URL} -c "DROP DATABASE IF EXISTS meal_planner;"

## db/create: create local db
.PHONY: db/create
db/create:
	psql ${PG_DATABASE_URL} -c "CREATE DATABASE meal_planner;"

## db/reset: reset the local db and setup fresh
.PHONY: db/reset
db/reset: db/drop db/create
	goose -dir ${SCHEMA_PATH} postgres "${DATABASE_URL}" up
	./scripts/populate_cuisines/run-local.sh

## migrate/%: goose migrate
.PHONY: migrate/%
migrate/%:
	goose -dir ${SCHEMA_PATH} postgres "${DATABASE_URL}" $(*)

## test: run all tests
.PHONY: test
test:
	./scripts/test_integration.sh

## build: build the application locally
.PHONY: build
build:
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## build/prod: build prod binary
.PHONY: build/prod
build/prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the application locally
.PHONY: run
run: build
	/tmp/bin/${BINARY_NAME}

## air: run the application with reloading on file changes
.PHONY: air
air:
	## Config is in .air.toml
	go run github.com/air-verse/air@latest
