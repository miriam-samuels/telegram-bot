build:
	@go build -o cmd/telegram_bot

run:
	@go run cmd/telegram_bot

test:
	@go test -v ./...