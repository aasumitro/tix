.PHONY: watch run critic lint tests api-spec mock migrate-down migrate-up migration-table build-fe

# Exporting bin folder to the path for makefile
export PATH   := $(PWD)/bin:$(PATH)
# Default Shell
export SHELL  := bash
# Type of OS: Linux or Darwin.
export OSTYPE := $(shell uname -s)

# --- Tooling & Variables ----------------------------------------------------------------
include ./misc/make/tools.Makefile

install-deps: gocritic gotestsum golangci-lint
deps: $(GOCRITIC) $(GOTESTSUM) $(GOLANGCI)
deps:
	@ echo "Required Tools Are Available"

migration-table:
	@if [ "$(name)" ]; then \
    	migrate create -ext sql -dir db/migrations $(name); \
   	fi

migrate-up:
	 migrate -database "${PQ_DSN}" -path db/migrations up

migrate-down:
	 migrate -database "${PQ_DSN}" -path db/migrations down

mock: $(MOCKERY)
	mockery --all --output=mocks/genmocks --outpkg=mocks

api-spec: tests
	@ echo "Re-generate API-Spec docs"
	@ swag init --parseDependency --parseInternal \
		--parseDepth 4 -g ./cmd/web/main.go

tests: $(GOTESTSUM) lint
	@ echo "Trying to run all tests cases"
	@ gotestsum --format standard-quiet \
		--hide-summary=skipped,output \
		-- -coverprofile=cover.out ./...
	@ rm cover.out

lint: $(GOLANGCI)
	@ echo "Trying apply linter"
	@ golangci-lint cache clean
	@ golangci-lint run -c .golangci.yaml ./...

critic: $(GOCRITIC)
	@ echo "Trying to critic the code"
	@ gocritic check -enableAll ./...

run: build-fe
	@echo "Run App"
	go mod tidy -compat=1.19
	go run ./cmd/web/main.go

build-fe:
	@ echo "Build Frontend"
	@ cd web && yarn install && yarn build

watch:
	air