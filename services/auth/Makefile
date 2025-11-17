.PHONY: build run fmt test migrate-up migrate-down

build:
	go build ./...

run:
	go run ./cmd/server

fmt:
	go fmt ./...

# test:
#     go test ./... -v

# migrate-up:
#     migrate -path ./migrations -database $(DATABASE_DSN) up

# migrate-down:
#     migrate -path ./migrations -database $(DATABASE_DSN) down
