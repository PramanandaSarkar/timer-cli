build:
	go build -o bin/timer-cli ./cmd/timer-cli

run:
	go run ./cmd/timer-cli

test:
	go test ./...

clean:
	rm -rf bin/