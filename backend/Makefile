.PHONY: up down run build test test-ci mocks clean db\:init db\:indexes db\:seed seed local-user

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

db\:init:
	@set -a && . ./.env && set +a && mongosh "$$MONGODB_DSN/$$MONGODB_DATABASE_NAME" < docs/mongodb/init.mongodb.js

db\:indexes:
	@set -a && . ./.env && set +a && mongosh "$$MONGODB_DSN/$$MONGODB_DATABASE_NAME" < docs/mongodb/indexes.mongodb.js

db\:seed:
	@set -a && . ./.env && set +a && mongosh "$$MONGODB_DSN/$$MONGODB_DATABASE_NAME" < docs/mongodb/seeds.mongodb.js

# Populates matches from frontend/public/data/worldcup2026.json.
# RESULTS=1 also loads the scores and marks the matches FINISHED.
seed:
	@set -a && . ./.env && set +a && go run ./scripts/seed $(if $(RESULTS),-results)

local-user:
	@test -n "$(CREDS)" || (echo "usage: make local-user CREDS=email:password" && exit 1)
	@set -a && . ./.env && set +a && go run scripts/create-local-user/main.go "$(CREDS)"
