BUILD_CONTEXT=./deploy
LINTER_CONCURRENCY ?= 4

.PHONY: go-inspect
go-inspect:
	@echo "Run static code analysis tools..."
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	GO111MODULE=on golangci-lint run \
				-j $(LINTER_CONCURRENCY) \
				--enable golint \
				--enable gosec \
				--enable gofmt \
				./...

.PHONY: go-test
go-test:
	@echo "Run all project tests..."
	GO111MODULE=on go test -v -p 1 ./...

.PHONY: go-build
go-build:
	@echo "Build project binaries..."
	GO112MODULE=on CGO_ENABLED=0 go build -v -o $(BUILD_CONTEXT)/medicare-api main.go

.PHONY: migrate
migrate:
	@echo "Migrate database..."
	GO111MODULE=on go run main.go -migrate

.PHONY: run
run:
	@echo "Run jobs..."
	GO111MODULE=on go run main.go

.PHONY: clean
clean:
	rm -rf $(BUILD_CONTEXT) 

