.PHONY: build run test lint docker clean

build:
	go build -o out/app cmd/app/main.go

run:
	ENV=local go run ./cmd/app

test:
	go test ./...

lint:
	golangci-lint run

docker:
	docker build -t insider-case .

clean:
	rm -rf out/

