APPLICATION_NAME := $(shell grep "const ApplicationName " version.go | sed -E 's/.*"(.+)"$$/\1/')
BIN_NAME=${APPLICATION_NAME}

BASE_VERSION := $(shell grep "const Version " version.go | sed -E 's/.*"(.+)"$$/\1/')

GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

default: help

get-deps: ## Install projects dependencies with Go Module
	@echo "Getting dependencies"
	go mod tidy
	go mod vendor

build: get-deps ## Build project for native production
	@echo "building ${BIN_NAME} ${BASE_VERSION}"
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY}" -o .bin/${BIN_NAME} ./

docker-build: build ## Build docker image
	sudo docker build -t ${APPLICATION_NAME}:${BASE_VERSION} ./

init: docker-build ## Initializing everything
	@echo "Initializing"
	sudo docker-compose up -d

test: ## Run project tests
	@echo "Running tests"
	mkdir -p ./.test/cover
	go test -race -coverpkg= ./... -coverprofile=./.test/cover/cover.out -v
