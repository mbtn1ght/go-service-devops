SHELL := /bin/sh

BRANCH ?= $(shell git branch --show-current)
VPS_HOST ?= 193.233.246.215
VPS_USER ?= root
VPS_DIR ?= /opt/service-devops
PNPM ?= corepack pnpm
GO_CACHE ?= /tmp/go-service-devops-go-build
SSH_FLAGS ?= -o BatchMode=yes -o ConnectTimeout=15
SSH = ssh $(SSH_FLAGS) $(VPS_USER)@$(VPS_HOST)

.DEFAULT_GOAL := help
.PHONY: help install frontend run up down restart logs ps test check-clean push update deploy

help:
	@printf '%s\n' \
	  'make run     Build the frontend and run the Go service locally.' \
	  'make up      Build and start the complete application with Docker Compose.' \
	  'make down    Stop local Docker Compose services.' \
	  'make logs    Follow local application logs.' \
	  'make test    Run Go tests.' \
	  'make push    Push the clean current branch to origin.' \
	  'make update  Update the VPS from the current branch and restart it.' \
	  'make deploy  Push, update the VPS, and verify the health check.'

install:
	cd web && $(PNPM) install --frozen-lockfile --allow-build=esbuild --allow-build=@clerk/shared

frontend: install
	cd web && $(PNPM) build

run: frontend
	GOCACHE=$(GO_CACHE) WEB_DIR=web/dist go run ./cmd

up:
	docker compose up -d --build --remove-orphans

down:
	docker compose down

restart: down up

logs:
	docker compose logs -f app

ps:
	docker compose ps

test:
	GOCACHE=$(GO_CACHE) go test ./...

check-clean:
	@test -z "$$(git status --porcelain)" || { echo 'Working tree is not clean. Commit or stash changes first.' >&2; exit 1; }

push: check-clean
	git push -u origin $(BRANCH)

update:
	$(SSH) 'set -eu; cd "$(VPS_DIR)"; git fetch --prune origin "$(BRANCH)"; git reset --hard "origin/$(BRANCH)"; git clean -fd; docker compose up -d --build --remove-orphans; docker compose ps; curl -fsS http://127.0.0.1:8080/healthz >/dev/null'

deploy: test push update
