.PHONY: frontend-server frontend-build frontend-deploy backend-build backend-test backend-deploy

## Frontend
frontend-server:
	$(MAKE) -C frontend server

frontend-build:
	$(MAKE) -C frontend build

frontend-deploy:
	$(MAKE) -C frontend deploy

## Backend
backend-build:
	$(MAKE) -C backend build

backend-test:
	$(MAKE) -C backend test

backend-deploy:
	$(MAKE) -C backend deploy
