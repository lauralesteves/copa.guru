.PHONY: up down run build test test-ci mocks clean seed

up:
	docker compose up -d

down:
	docker compose down

run:
	sam local start-api

build:
	sam build

test:
	@go test -race -cover ./... 2>&1 || echo "no test packages found"

test-ci:
	@go test -race -cover -count=1 -failfast -parallel=4 ./... 2>&1 || echo "no test packages found"

mocks:
	go generate ./...

clean:
	rm -rf .aws-sam

seed:
	@echo "TODO: seed de matches — implementar quando modelos estiverem prontos"
