build:
	@go build -o bin/toronto-bookings

run: build
	@./bin/toronto-bookings

test:
	@go test -v ./...