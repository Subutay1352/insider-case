.PHONY: build run test lint docker clean swagger

build:
	go build -o out/app cmd/app/main.go

run:
	ENV=local go run ./cmd/app

test:
	go test ./...

lint:
	golangci-lint run

swagger:
	~/go/bin/swag init -g cmd/app/main.go -o docs

docker:
	docker build -t insider-case .

clean:
	rm -rf out/

