.PHONY: run test test-race vet fmt-check

run:
	go run ./cmd/api

test:
	go test ./...

test-race:
	go test -race ./...

vet:
	go vet ./...

fmt-check:
	@test -z "$$(gofmt -l .)"
