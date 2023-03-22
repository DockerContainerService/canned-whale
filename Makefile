NAME = canned-whale

OS = linux
ARCHS = amd64 arm64

.DEFAULT_GOAL := help

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

all: release

build: deps ## Build the project
	go build -ldflags "-s -w"

release-amd64: clean deps ## Generate release for linux amd64
	mkdir -p build;
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o build/$(NAME)-linux-amd64;
	upx -9 build/$(NAME)-linux-amd64;

release: clean deps ## Generate releases for unix systems
	@for arch in $(ARCHS);\
	do \
		for os in $(OS);\
		do \
			echo "Building $$os-$$arch"; \
			mkdir -p build; \
			GOOS=$$os GOARCH=$$arch go build -ldflags "-s -w" -o build/$(NAME)-$$os-$$arch; \
			upx -9 build/$(NAME)-$$os-$$arch; \
		done \
	done

test: deps ## Execute tests
	go test ./...

deps: ## Install dependencies using go get
	go get -d -v -t ./...

clean: ## Remove building artifacts
	rm -rf build
	rm -f $(NAME)
