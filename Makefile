.PHONY: init \
       frontend-server frontend-build frontend-deploy \
       backend-up backend-down backend-run backend-build backend-test backend-test-ci \
       backend-mocks backend-clean backend-db-init backend-local-user backend-deploy

## Setup
init:
	cp scripts/pre-push .git/hooks/pre-push
	chmod +x .git/hooks/pre-push
	@echo "Git hooks instalados."

## Frontend
frontend-server:
	$(MAKE) -C frontend server

frontend-build:
	$(MAKE) -C frontend build

frontend-deploy:
	$(MAKE) -C frontend deploy

## Backend
backend-up:
	$(MAKE) -C backend up

backend-down:
	$(MAKE) -C backend down

backend-run:
	$(MAKE) -C backend run

backend-build:
	$(MAKE) -C backend build

backend-test:
	$(MAKE) -C backend test

backend-test-ci:
	$(MAKE) -C backend test-ci

backend-mocks:
	$(MAKE) -C backend mocks

backend-clean:
	$(MAKE) -C backend clean

backend-db-init:
	$(MAKE) -C backend db:init
	$(MAKE) -C backend db:indexes
	$(MAKE) -C backend db:seed

backend-local-user:
	$(MAKE) -C backend local-user CREDS="$(CREDS)"

backend-deploy:
	$(MAKE) -C backend deploy
