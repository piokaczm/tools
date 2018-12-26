VERSION	?= $(shell date +"%Y-%m-%d")
BUILD		?= dev
HASH		?= $(shell git rev-parse HEAD)
PROJECT = foureyes

LDFLAGS = -X main.buildDate=$(VERSION)	\
	  -X main.release=$(BUILD)	\
	  -X main.commitSha=$(HASH)

test: ## Run all tests
	GO111MODULE=on go test -cover ./...

run: build-binary ## Run the project locally
	./$(PROJECT)

build: ## Build binary with compile time flags
	GO111MODULE=on CGO_ENABLED=0 go build -o $(PROJECT) cmd/$(PROJECT)/*.go