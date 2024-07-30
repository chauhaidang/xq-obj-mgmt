build:
	go build -o cmd/bin/app

run: build
	cmd/bin/app

test:
	go test -v ./...