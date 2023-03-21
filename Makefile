NAME = canned-whale

OS = linux
ARCHS = 386 arm amd64 arm64
ARCHSW = 386 amd64

.DEFAULT_GOAL := help

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

all: build release release-windows

build: deps ## Build the project
	go build -ldflags "-s -w"

release: clean deps ## Generate releases for unix systems
	@for arch in $(ARCHS);\
	do \
		for os in $(OS);\
		do \
			echo "Building $$os-$$arch"; \
			mkdir -p build; \
			GOOS=$$os GOARCH=$$arch go build -ldflags "-s -w" -o build/$(NAME)-$$os-$$arch; \
		done \
	done

release-windows: clean deps ## Generate release for windows
	@for arch in $(ARCHSW);\
	do \
		echo "Building windows-$$arch"; \
		mkdir -p build; \
		GOOS=windows GOARCH=$$arch go build -ldflags "-s -w" -o build/$(NAME)-windows-$$arch.exe; \
	done

test: deps ## Execute tests
	go test ./...

deps: ## Install dependencies using go get
	go get -d -v -t ./...

clean: ## Remove building artifacts
	rm -rf build
	rm -f $(NAME)
