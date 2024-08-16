# Change these variables as necessary.
MAIN_PACKAGE_PATH := .
BINARY_NAME := mealorg_server
SCHEMA_PATH := sql/schema
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
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0 generate

## db/drop: drop local db
.PHONY: db/drop
db/drop:
	psql -h localhost -d postgres -c "DROP DATABASE IF EXISTS ${DB_NAME};"

## db/create: create local db
.PHONY: db/create
db/create:
	psql -h localhost -d postgres -c "CREATE DATABASE ${DB_NAME};"

## db/reset: reset the local db and setup fresh
.PHONY: db/reset
db/reset: db/drop db/create
	./scripts/goose.sh up
	./scripts/populate_cuisines/run-local.sh

## migrate/%: goose migrate
.PHONY: migrate/%
migrate/%:
	./scripts/goose.sh $(*)

## test: run all tests
.PHONY: test
test:
	./scripts/test_integration.sh

## templ: generate templ code
.PHONY: templ
templ:
	go run github.com/a-h/templ/cmd/templ@latest generate

## tailwind: compile tailwind css
.PHONY: tailwind
tailwind:
	npx tailwindcss -i ./assets/css/input.css -o ./assets/css/styles.css --minify

## build: build the application locally
.PHONY: build
build: sqlc templ tailwind
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## build/prod: build prod binary
.PHONY: build/prod
build/prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the application locally
.PHONY: run
run: build
	/tmp/bin/${BINARY_NAME}

## live/templ: run templ generation in watch mode to detect all .templ changes
.PHONY: live/templ
live/templ:
	go run github.com/a-h/templ/cmd/templ@latest generate --watch --proxy="http://localhost:8080" --open-browser=false

## live/server: run the application with reloading on file changes
.PHONY: live/server
live/server:
	## Config is in .air.toml
	go run github.com/air-verse/air@v1.52.3 \
  --build.cmd "go build -o tmp/bin/${BINARY_NAME} && templ generate --notify-proxy" \
	--build.bin "tmp/bin/${BINARY_NAME}" \
	--build.delay "100" \
  --build.exclude_dir "node_modules,sql,scripts,tests" \
  --build.include_ext "go" \
  --build.exclude_regex "" \
  --build.stop_on_error "false" \
  --build.post_cmd "pkill ${BINARY_NAME}" \
  --misc.clean_on_exit true

## live/tailwind: run the application with reloading on file changes
.PHONY: live/tailwind
live/tailwind:
	npx tailwindcss -i ./assets/css/input.css -o ./assets/css/styles.css --minify --watch

## live/assets: watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
.PHONY: live/assets
live/assets:
	# Perhaps not necessary unless I have a separate js/css compilation process
	go run github.com/air-verse/air@v1.52.3 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "assets" \
	--build.include_ext "js,css"

## live: start all watch processes in parallel
.PHONY: live
live:
	make -j3 live/templ live/server live/tailwind
