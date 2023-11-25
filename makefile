.PHONY: help
LOCALBIN=$(shell pwd)/bin
BINARY=simplerestserver

version:
	@go version

fmt: version
	@test -s $(LOCALBIN)/gofumpt || GOBIN=$(LOCALBIN) go install mvdan.cc/gofumpt@latest
	$(LOCALBIN)/gofumpt -l -w .

lint: fmt
	@test -s $(LOCALBIN)/golangci-lint  || \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
			| sh -s -- -b $(LOCALBIN) v1.55.2
	$(LOCALBIN)/golangci-lint run ./...

vulncheck: lint
	@test -s $(LOCALBIN)/govulncheck || \
		GOBIN=$(LOCALBIN) go install golang.org/x/vuln/cmd/govulncheck@latest
	$(LOCALBIN)/govulncheck ./...

check:	vulncheck
	go generate ./...
	go vet ./...
	go clean ./...
	go mod tidy
	go mod download

build: check
	CGO_ENABLED=0 go build -ldflags='-s -w' -o $(LOCALBIN)/$(BINARY) .
	@echo "build completed. and binary place at  ::"$(LOCALBIN)/$(BINARY) 

run: build
	@$(LOCALBIN)/$(BINARY)