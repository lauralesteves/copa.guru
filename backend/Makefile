.PHONY: up down run build test test-ci mocks clean seed local-user

up:
	docker compose up -d

down:
	docker compose down

run: build
	sam local start-api

build: clean
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

local-user:
	@test -n "$(CREDS)" || (echo "usage: make local-user CREDS=email:password" && exit 1)
	@set -a && . ./.env && set +a && go run scripts/create-local-user/main.go "$(CREDS)"
