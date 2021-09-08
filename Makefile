CMD := locker

build:
	go build -o $(CMD) ./cmd/main.go

clean:
	rm $(CMD)

lint:
	golangci-lint run ./...