build:
	go build -o out/app ./cmd/app

run:
	ENV=local go run ./cmd/app

test:
	go test ./...

lint:
	~/go/bin/golangci-lint run

swagger:
	~/go/bin/swag init -g cmd/app/main.go -o docs

docker:
	docker build -t insider-case .

clean:
	rm -rf out/

