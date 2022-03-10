lint:
	golangci-lint run --timeout=2m0s

.PHONY: test
test: lint
	go test -coverprofile cover.out ./...
	go tool cover -func cover.out
